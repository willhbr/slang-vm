require_relative './objects'
require_relative './builtins'

class Resolver
  def initialize(defs=Hash.new)
    @next_id = 0
    @current_module = nil
  end

  def resolve_top_level(ast)
    resolve(ast, [Hash.new], true)
  end

  private

  def resolve(ast, scopes=[Hash.new], top_level=false)
    case ast
    when Vector
      ast.map { |node| resolve(node, scopes) }
    when Array
      resolve_call(ast, scopes, top_level)
    when Hash
      res = Hash.new
      ast.each do |k, v|
        res[resolve(k, scopes)] = resolve(v, scopes)
      end
      res
    when Identifier
      return if Identifier::KEYWORDS.include? ast.value
      return if bind_to_local(ast, scopes)
      if iden = Defs.get?(ast, @current_module.value)
        ast.code = iden.code
        ast.to_module_only!
        return
      end
      raise "Undefined var: #{ast.value}"
    else
      ast
    end
  end

  def resolve_fn(ast, scopes)
    scopes = rescope(scopes)
    args = ast[1]
    raise "Args must be vector, not #{args.class}" unless args.is_a? Vector
    args.each do |arg|
      raise "Arg must be identifier: #{arg}" unless arg.is_a? Identifier
      set_binding(arg, scopes)
    end
    ast[2..-1].each do |expr|
      resolve(expr, scopes)
    end
  end

  def resolve_call(ast, scopes, top_level=false)
    first = ast.first
    return unless first
    if first.is_a? Identifier
      case first.value
      when 'let'
        binds = ast[1]
        scopes = rescope(scopes)
        binds.each_slice(2) do |slice|
          name, expr = slice
          resolve(expr, scopes)
          raise "Not an identifier: #{name}" unless name.is_a? Identifier
          set_binding(name, scopes)
        end
        ast[2..-1].each do |node|
          resolve(node, scopes)
        end
      when 'module'
        raise "Can only define module at top-level" unless top_level
        first, name = ast
        name.to_module_only!
        name.code = Defs.define_module(name)
        @current_module = name
      when 'def'
        raise "Can only def at top-level" unless top_level
        first, name, value = ast
        resolve(value, scopes)
        raise 'Cannot define outside of module' unless @current_module
        raise "Cannot define in other module: #{name.mod}" if name.module
        if iden = Defs.get?(name, @current_module.value)
          raise "Already defined #{name.value} on @{iden.location}"
        end
        name.code = Defs.set_and_return(name, @current_module.value)
      when 'fn'
        resolve_fn(ast, scopes)
      when 'do'
        ast[1..-1].each { |node| resolve_fn(node, scopes, top_level) }
      else
        ast.each { |node| resolve(node, scopes) }
      end
    else
      ast.each { |node| resolve(node, scopes) }
    end
  end

  def bind_to_local(identifier, scopes)
    scopes.reverse.each do |scope|
      if bound = scope[identifier.value]
        identifier.code = bound.code
        return true
      end
    end
    return false
  end

  def set_binding(identifier, scopes)
    @next_id += 1
    identifier.code = @next_id
    scopes[-1][identifier.value] = identifier
  end

  def rescope(scopes)
    scopes + [Hash.new]
  end
end

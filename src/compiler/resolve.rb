require_relative './objects'
require_relative './builtins'

class ResolverState
  attr_accessor :scopes
  def initialize
    @scopes = [ Hash.new ]
    @top_level = true
  end

  def rescope(&block)
    @scopes.push Hash.new
    result = block.()
    @scopes.pop
    result
  end
end

class Resolver
  def initialize(defs=Hash.new)
    @next_id = 0
    @current_module = nil
    @state = ResolverState.new
  end

  def resolve_top_level(ast)
    resolve(ast, true)
  end

  private

  def resolve(ast, top_level=false)
    case ast
    when Vector
      ast.map { |node| resolve(node) }
    when Array
      resolve_call(ast, top_level)
    when Hash
      res = Hash.new
      ast.each do |k, v|
        res[resolve(k)] = resolve(v)
      end
      res
    when Identifier
      return if Identifier::KEYWORDS.include? ast.value
      return if bind_to_local(ast)
      raise "No module defined!" unless @current_module
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

  def resolve_fn(ast)
    @state.rescope do
      args = ast[1]
      raise "Args must be vector, not #{args.class}" unless args.is_a? Vector
      args.each do |arg|
        raise "Arg must be identifier: #{arg}" unless arg.is_a? Identifier
        set_binding(arg)
      end
      ast[2..-1].each do |expr|
        resolve(expr)
      end
    end
  end

  def resolve_call(ast, top_level=false)
    first = ast.first
    return unless first
    if first.is_a? Identifier
      case first.value
      when 'let'
        binds = ast[1]
        @state.rescope do
          binds.each_slice(2) do |slice|
            name, expr = slice
            resolve(expr)
            raise "Not an identifier: #{name}" unless name.is_a? Identifier
            set_binding(name)
          end
          ast[2..-1].each do |node|
            resolve(node)
          end
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
        resolve(value)
        raise 'Cannot define outside of module' unless @current_module
        raise "Cannot define in other module: #{name.mod}" if name.module
        if iden = Defs.get?(name, @current_module.value)
          raise "Already defined #{name.value} on #{iden.location}"
        end
        name.code = Defs.set_and_return(name, @current_module.value)
      when 'fn'
        resolve_fn(ast)
      when 'do'
        ast[1..-1].each { |node| resolve_fn(node) }
      else
        ast.each { |node| resolve(node) }
      end
    else
      ast.each { |node| resolve(node) }
    end
  end

  def bind_to_local(identifier)
    @state.scopes.reverse.each do |scope|
      if bound = scope[identifier.value]
        identifier.code = bound.code
        return true
      end
    end
    return false
  end

  def set_binding(identifier)
    @next_id += 1
    identifier.code = @next_id
    @state.scopes[-1][identifier.value] = identifier
  end
end

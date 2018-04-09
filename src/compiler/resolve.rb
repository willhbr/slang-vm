require_relative './objects'
require_relative './builtins'

class ResolverState
  attr_accessor :scopes
  attr_accessor :in_func

  def initialize
    @scopes = [ Hash.new ]
    @top_level = true
    @in_func = 0
  end

  def rescope(&block)
    @scopes.push Hash.new
    result = block.()
    @scopes.pop
    result
  end
end

class Resolver
  def initialize
    @next_id = 0
    @current_module = nil
    @on_local_bind = []
    @state = ResolverState.new
  end

  def resolve_top_level(ast)
    @next_id = 0
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
      return if Identifier::KEYWORDS.include? ast.whole
      return if bind_to_local(ast)
      raise "No module defined!" unless @current_module
      return if Defs.get_module?(@current_module, ast)
      return if Defs.get_module_def?(@current_module, ast)
      raise "Undefined var: #{ast.whole} #{ast.location}"
    else
      ast
    end
  end

  def resolve_fn(ast)
    @state.in_func += 1
    @state.rescope do
      args = ast[1]
      outside_code = @next_id
      raise "Args must be vector, not #{args.class} #{ast[0].location}" unless args.is_a? Vector
      args.each do |arg|
        raise "Arg must be identifier: #{arg}" unless arg.is_a? Identifier
        set_binding(arg)
      end
      captured = []
      # This is a hackety hack
      @on_local_bind << lambda do |iden|
        if iden.code < outside_code
          captured << iden
        end
      end
      args.push(ClosureArgs.new captured)
      ast[2..-1].each do |expr|
        resolve(expr)
      end
      @on_local_bind.pop
    end
    @state.in_func -= 1
  end

  def resolve_call(ast, top_level=false)
    first = ast.first
    return unless first
    if first.is_a? Identifier
      case first.whole
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
      when 'alias'
        raise 'cannot alias outside module' unless @current_module
        mod = ast[1]
        raise 'Can only alias modules' unless mod.is_a? Identifier
        resolve(mod)
        alias_to = ast[2]
        raise 'Can only alias modules' unless alias_to.is_a?(Identifier)
        Defs.alias(@current_module, mod, alias_to)
      when 'import'
        raise 'cannot import outside module' unless @current_module
        mod = ast[1]
        raise 'Can only alias modules' unless mod.is_a? Identifier
        resolve(mod)
        Defs.import(@current_module, mod)
      when 'module'
        raise "Can only define module at top-level" unless top_level
        _, name = ast
        Defs.define_module(name)
        @current_module = name
      when 'def'
        raise "Can only def at top-level" unless top_level
        first, name, value = ast
        resolve(value)
        raise 'Cannot define outside of module' unless @current_module
        raise "Cannot define in other module: #{name.mod}" unless name.no_module?
        if Defs.get_module_def?(@current_module, name)
          raise "Already defined #{name.whole} on #{iden.location}"
        end
        Defs.def_def(@current_module, name)
      when 'recur'
        if @state.in_func != 0
          ast[1..-1].each { |node| resolve(node, top_level) }
        else
          raise "Can't recur outside function"
        end
      when 'fn'
        resolve_fn(ast)
      when 'do'
        ast[1..-1].each { |node| resolve(node, top_level) }
      else
        if Identifier::KEYWORDS.include? first.whole
          ast[1..-1].each { |node| resolve(node) }
        else
          ast.each { |node| resolve(node) }
        end
      end
    else
      ast.each { |node| resolve(node) }
    end
  end

  def bind_to_local(identifier)
    @state.scopes.reverse.each do |scope|
      if bound = scope[identifier.whole]
        identifier.code = bound.code
        @on_local_bind.each do |func|
          func.(identifier)
        end
        return true
      end
    end
    return false
  end

  def set_binding(identifier)
    identifier.code = @next_id
    @next_id += 1
    @state.scopes[-1][identifier.whole] = identifier
  end
end

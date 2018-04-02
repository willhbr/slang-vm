require_relative './objects'

class Resolver
  KEYWORDS = ['if']
  def initialize(defs=Hash.new)
    @defs = defs
    @next_id = 0
    @defs.each do |_, iden|
      @next_id += 1
      iden.code = @next_id
    end
  end

  def resolve(ast, scopes=[Hash.new])
    case ast
    when Vector
      ast.map { |node| resolve(node, scopes) }
    when Array
      resolve_call(ast, scopes)
    when Hash
      res = Hash.new
      ast.each do |k, v|
        res[resolve(k, scopes)] = resolve(v, scopes)
      end
      res
    when Identifier
      if KEYWORDS.include? ast.value
        return
      end
      unless bind_to_previous(ast, scopes)
        raise "Undefined var: #{ast}"
      end
    else
      ast
    end
  end

  def resolve_call(ast, scopes)
    first = ast.first
    return unless first
    if first.is_a?(Identifier) && first.value == 'let'
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
    elsif first.is_a?(Identifier) && first.value == 'def'
      first, name, value = ast
      resolve(value, scopes)
      set_binding(name, [@defs])
    else
      ast.map { |node| resolve(node, scopes) }
    end
  end

  def bind_to_previous(identifier, scopes)
    scopes.reverse.each do |scope|
      if bound = scope[identifier.value]
        identifier.code = bound.code
        return true
      end
    end
    if bound = @defs[identifier.value]
      identifier.code = bound.code
      return true
    end
    raise "Undefined variable: #{identifier.inspect}"
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

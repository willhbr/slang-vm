require_relative './objects'
require_relative './builtins'

class MacroExpander
  MACROS = {
    defn: lambda do |ast|
      d, name, args = ast
      body = ast[3..-1]
      [Identifier.from(d, "def"), name,
        [Identifier.from(d, "fn"), args] + body]
    end,
    "->>": lambda do |ast|
      ast[1..-1].reduce do |inner, call|
        if call.is_a? Array
          call + [inner]
        else
          [call, inner]
        end
      end
    end,
    "->": lambda do |ast|
      ast[1..-1].reduce do |inner, call|
        if call.is_a? Array
          [call[0], inner] + call[1..-1]
        else
          [call, inner]
        end
      end
    end,
    "unless": lambda do |ast|
      unless_, cond, then_, other = ast
      [Identifier.from(unless_, 'if'), cond, other || Identifier.new('nil', nil), then_]
    end
  }

  def initialize
    @macros = Hash.new
    MACROS.each do |name, macro|
      @macros[name.to_s] = macro
    end
  end

  def expand_top_level(ast)
    expand(ast)
  end

  private

  def expand(ast)
    case ast
    when Vector
      Vector.new ast.map { |node| expand(node) }
    when Array
      expand_call(ast)
    when Hash
      res = Hash.new
      ast.each do |k, v|
        res[expand(k)] = expand(v)
      end
      res
    when Identifier
      ast
    else
      ast
    end
  end

  def expand_call(ast)
    first = ast.first
    return unless first
    if first.is_a? Identifier
      if macro = @macros[first.whole]
        return macro.(ast)
      end
    end
    ast.map { |node| expand(node) }
  end
end

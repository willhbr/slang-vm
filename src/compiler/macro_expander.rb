require_relative './objects'
require_relative './builtins'
require_relative './ast_processor'

class MacroExpander
  include ASTProcessor

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
    end,
    "deftype": lambda do |ast|
      deftype, name = ast
      attrs = ast[2..-1]
      args = (0...attrs.size).map { |a| kw(:"arg_#{a}") }
      [kw(:do),
        ast,
        [kw(:def), from(name, :"#{name.whole}.new"),
          [kw(:fn), vec(*args),
            [kw(:"new-instance"), name] + args]]]
    end
  }

  def initialize
    @macros = Hash.new
    MACROS.each do |name, macro|
      @macros[name.to_s] = macro
    end
  end

  def process_array(ast, top_level)
    first = ast.first
    return unless first
    if first.is_a? Identifier
      if macro = @macros[first.whole]
        return macro.(ast)
      end
    end
    ast.map { |node| process(node) }
  end

  def process_other(ast, top_level)
    passthrough_nested(ast)
  end

  def self.kw(name)
    Identifier.new name.to_s, nil
  end

  def self.vec(*items)
    Vector.new items
  end

  def self.from(iden, name)
    Identifier.from(iden, name.to_s)
  end
end

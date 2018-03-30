require_relative './resolve'
require_relative './objects'

class ExpectedEOF < Exception
end

class Scanner
  class Token
    attr_reader :type
    attr_reader :value

    def initialize(type, value)
      @type = type
      @value = value
    end

    def to_s
      "#{@type} #{@value}"
    end

    def inspect; to_s end
  end

  def initialize(filename, contents)
    @filename = filename
    @contents = contents
    @index = 0
  end

  def read
    tokens = []
    while @index < @contents.size
      begin
        if token = next_token
          tokens.push token
        end
      rescue Exception => e
        break
      end
    end
    tokens << sym(:EOF)
    tokens
  end

  def next_token
    c = advance!
    case c
    when '('
      return sym(:'(')
    when ')'
      return sym(:')')
    when ')'
      return sym(:")")
    when '['
      return sym(:"[")
    when ']'
      return sym(:"]")
    when '\''
      return sym(:"'")
    when '`'
      return sym(:"`")
    when '@'
      return sym(:"@")
    when '~'
      if peek? == '@'
        advance?
        return sym(:"~@")
      else
        return sym(:"~")
      end
    when '{'
      return sym(:"{")
    when '}'
      return sym(:"}")
    when ':'
      iden = identifier("")
      return sym(:ATOM, iden)
    when ' ', "\t", "\n", ','
      return nil
    when '"'
      return sym(:STRING, string)
    when '#'
      return sym(:READER_MACRO)
    when '/'
      return sym(:REGEX_LITERAL, regex)
    when ';'
      comment
      return nil
    else
      if /^[0-9]$/ =~ c
        return sym(:NUMBER, number)
      elsif is_iden_start(c)
        val = identifier
        if val.end_with? ':'
          return sym(:KW_ARG, val[0..-2])
        else
          return sym(:IDENTIFIER, val)
        end
      else
        raise "Invalid character! '#{c}'"
      end
    end
  end

  def string
    start = @index - 1
    advance! while peek? != '"'
    str = @contents[start..@index]
    advance!
    str
  end

  def number
    start = @index - 1
    advance! while peek? =~ /\d/
    str = @contents[start..@index]
    str.to_i
  end

  def identifier
    start = @index - 1
    advance! while is_iden_start(peek?)
    return @contents[start..@index - 1]
  end

  def is_iden_start(char)
    'qwertyuiopasdfghjklzxcvbnm!@#$%^&*-_=+\|:?/,<>.'.include? char
  end

  def sym(type, value=nil)
    Token.new(type, value)
  end

  def peek?
    @contents[@index]
  end

  def lookahead
    @contents[@index + 1]
  end

  def advance!
    c = @contents[@index]
    raise "EOF" unless c
    @index += 1
    c
  end
end

class Parser
  def initialize(tokens)
    @tokens = tokens
    @index = 0
  end

  def program
    prog = []
    loop do
      begin
        prog << object
      rescue ExpectedEOF
        break
      end
    end
    prog
  end

  def object
    symbol = pop_sym?
    case symbol.type
      when :'('
        return list(:')', Array.new)
      when :'['
        return list(:']', Vector.new)
      when :'{'
        return map
      when :"'", :'`'
        o = object
        raise "EOF" unless o
        [:quote, object]
      when :'~'
        o = object
        raise "EOF" unless o
        [:unquote, object]
      when :'~@'
        o = object
        raise "EOF" unless o
        [:'unquote-splice', object]
      when :READER_MACRO
        reader_macro
      when :IDENTIFIER
        identifier(symbol)
      when :NUMBER
        number(symbol)
      when :STRING
        string(symbol)
      when :ATOM
        atom(symbol)
      when :KW_ARG
        kw_arg(symbol)
      when :EOF
        raise ExpectedEOF.new
      else
        raise "Syntax error, oops! #{symbol.type} #{symbol.value}"
    end
  end

  def identifier(symbol)
    Identifier.new symbol.value
  end

  def map
    into = Hash.new
    loop do
      break if !(sym = peek) || sym.type == :'}'
      raise 'EOF' unless peek
      k = object
      raise 'EOF' unless peek
      raise 'unexpected }' if !(sym = peek) || sym.type == :'}'
      v = object
      into[k] = v
    end
    pop_sym?
    into
  end

  def list(terminator, into)
    while (sym = peek) && sym.type != terminator
      into << object
    end
    raise 'EOF' if peek.nil?
    pop_sym?
    into
  end

  def number(token)
    token.value.to_i
  end

  def string(token)
    token.value
  end

  def kw_arg(token)
    Atom.new token.value, true
  end

  def atom(token)
    Atom.new token.value
  end

  def peek
    @tokens[@index]
  end

  def pop_sym?
    sym = @tokens[@index]
    @index += 1
    sym
  end
end

filename = ARGV[0]
scanner = Scanner.new filename, File.read(filename)
tokens = scanner.read
tree = Parser.new(tokens).program

res = Resolver.new({
  "print" => Identifier.new("print")
})
tree.map do |node|
  res.resolve(node)
end

puts tree.map(&:inspect).join "\n"

require_relative './objects'
require_relative './program'

class ExpectedEOF < Exception
end

class Location
  def initialize(filename, line, column)
    @filename = filename
    @line = line
    @column = column
  end

  def to_s
    "#{@filename}@#{@line}:#{@column}"
  end
end

class Scanner
  class Token
    attr_reader :type
    attr_reader :value
    attr_reader :location

    def initialize(type, value, location)
      @type = type
      @value = value
      @location = location
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
    @column = 1
    @line = 1
  end

  def read
    tokens = []
    while @index < @contents.size
      begin
        if token = next_token
          tokens.push token
        end
      rescue ExpectedEOF
        break
      end
    end
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
      iden = atom
      return sym(:ATOM, iden)
    when ' ', "\t", ','
      return nil
    when "\n"
      @column = 1
      @line += 1
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
    start = @index
    advance! while peek? != '"'
    str = @contents[start..@index - 1]
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
    advance! while is_iden_middle(peek?)
    return @contents[start..@index - 1]
  end

  def atom
    start = @index
    advance! while is_iden_middle(peek?)
    return @contents[start..@index - 1]
  end

  def comment
    advance! while peek? != "\n"
    @column = 1
    @line += 1
  end

  def is_iden_middle(char)
    '1234567890'.include?(char) || is_iden_start(char)
  end

  def is_iden_start(char)
    'QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm!@#$%^&*-_=+\|:?/,<>.'.include? char
  end

  def sym(type, value=nil)
    Token.new(type, value, Location.new(@filename, @line, @column))
  end

  def peek?
    @contents[@index]
  end

  def lookahead
    @contents[@index + 1]
  end

  def advance!
    c = @contents[@index]
    raise "Unexpected EOF" unless c
    @column += 1
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
    while o = object
      prog << o
    end
    prog
  end

  def object
    symbol = pop_sym?
    return nil unless symbol
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
      else
        raise "Syntax error: Unexpected #{symbol.type} '#{symbol.value}' (#{symbol.location})"
    end
  end

  def identifier(symbol)
    Identifier.new symbol.value, symbol.location
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



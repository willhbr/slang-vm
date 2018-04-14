require_relative './builtins'

class Program
  attr_reader :position

  def initialize
    @strings = Hash.new # string to const id
    @string_code = -1
    @buffer = Array.new
    @position = 0
  end

  def add_string(string)
    if code = @strings[string]
      code
    else
      @string_code += 1
      @strings[string] = @string_code
      @string_code
    end
  end

  def <<(code)
    @buffer.push code
    @position += code.size
  end

  def bytes
    buff = []
    @buffer.each do |code|
      code >> buff
    end
    buff
  end

  def string_bytes
    if @strings.size != @string_code + 1
      raise "Strings don't match code length"
    end
    buff = []
    buff << @strings.size

    strings = Array.new @strings.size
    @strings.each do |string, idx|
      strings[idx] = string
    end

    strings.each do |string|
      size = string.bytesize
      buff << size
      string.each_byte do |byte|
        buff << byte
      end
    end
    buff
  end

  def write_to(file)
    strings = self.string_bytes
    bytes = self.bytes
    file.write(strings.pack('C' * strings.size))
    file.write([ Defs.def_count ].pack('C'))
    file.write(bytes.pack('C' * bytes.size))
  end

  def print
    pos = 0
    @buffer.each do |code|
      puts code.to_s pos
      pos += code.size
    end
  end
end

class Arg
  def self.local(num, debug)
    if num < 0 || num > 255
      raise "Local arg code exceeds byte size: #{num} #{debug}"
    end
    new(:local, num, debug)
  end

  def self.global(num, debug)
    # TODO increase to 2-4 bytes
    if num < 0 || num > 255
      raise "Global def code exceeds byte size: #{num} #{debug}"
    end
    new(:global, num, debug)
  end

  def self.string(num, debug)
    # TODO increase to 2-4 bytes
    if num < 0 || num > 255
      raise "String code exceeds byte size: #{num} #{debug}"
    end
    new(:string, num, debug)
  end

  def self.atom(num, debug)
    # TODO increase to 2-4 bytes
    if num < 0 || num > 255
      raise "Atom code exceeds byte size: #{num} #{debug}"
    end
    new(:atom, num, debug)
  end

  def self.integer(value)
    new(:integer, value, value.to_s)
  end

  def initialize(type, code, debug)
    @type = type
    @code = code
    @debug = debug
  end

  def bytes_into(buffer)
    case @type
    when :integer
      if @code > 255 || @code < 0
        [ @code ].pack('Q>').split('').each do |byte|
          buffer << byte.ord
        end
      else
        buffer << @code
      end
    else
      buffer << @code
    end
  end

  def to_s
    "#{@code} (#{@debug})"
  end
end

class Code
  LOAD_LOCAL = 1
  LOAD_DEF = 2
  STORE = 3
  INVOKE = 4
  CONST_I = 5
  CONST_I_BIG = 6
  CONST_S = 7
  CONST_A = 8
  CONST_TRUE = 9
  CONST_FALSE = 10
  CONST_NIL = 11
  JUMP = 12
  JUMP_BACK = 13
  AND = 14
  RETURN = 15
  CLOSURE = 16
  PROTOCOL_CLOSURE = 17
  NEW_MAP = 18
  NEW_VECTOR = 19
  NEW_LIST = 20
  DEFINE = 21
  TYPE = 22
  INSTANCE = 23
  IMPLEMENT = 24
  RAISE = 25
  TRY = 26
  END_TRY = 27
  DISCARD = 28

  CODE_VALUES = Hash.new
  CODE_NAMES = Hash.new

  class << self
    Code.constants(false).each do |const|
      next if const == :CODE_VALUES || const == :CODE_NAMES
      val = Code.const_get const
      Code::CODE_VALUES[const] = val
      Code::CODE_NAMES[val] = const
      define_method(const) do |*args|
        Code.new val, *args
      end
    end
  end

  attr_reader :args
  attr_reader :debug
  def initialize(code, args=[], debug=nil)
    @code = code
    @args = args.is_a?(Array) ? args : [args]
    if @args.any?(&:nil?)
      raise "No args can be nil: #{args.inspect} #{debug.inspect}"
    end
    if debug
      @debug = debug.is_a?(Array) ? debug : [debug]
    else
      @debug = nil
    end
  end

  def to_s(pos='')
    name = self.class.stringify(@code)
    if @debug
      args = @args.zip(@debug).map { |a, d| "#{a}: #{d || '?'}" }.join("\t")
    else
      args = @args.map(&:to_s).join("\t")
    end
    '%4s %11s %s' % [pos.to_s, name, args]
  end

  def self.stringify(code)
    name = CODE_NAMES[code]
    unless name
      raise "Invalid instruction: #{code}"
    end
    name
  end

  def >>(buffer)
    buffer << @code
    @args.each do |arg|
      if arg.is_a?(Arg)
        arg.bytes_into(buffer)
      else
        buffer << arg
      end
    end
  end

  def size
    @args.size + 1
  end
end

if __FILE__==$0
  File.open(ARGV[0], 'w') do |file|
    file.print "package op_codes\nconst ("
    Code::CODE_VALUES.each do |name, value|
      file.puts "#{name} = #{value}"
    end
    file.puts ')'

    file.puts 'func ToString(code byte) string {'
    file.puts 'switch code {'
    Code::CODE_VALUES.each do |name, value|
      file.puts "case #{name}: return \"#{name}\""
    end
    file.puts 'default: return "UNKNOWN"'
    file.puts '}}'
  end
end

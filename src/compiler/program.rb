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


class Code
  LOAD_LOCAL = 1
  LOAD_DEF = 2
  STORE = 3
  INVOKE = 4
  APPLY = 5
  CONST_I = 6
  CONST_S = 7
  CONST_A = 8
  CONST_TRUE = 9
  CONST_FALSE = 10
  CONST_NIL = 11
  JUMP = 12
  AND = 13
  OR = 14
  RETURN = 15
  CLOSURE = 16
  NEW_MAP = 17
  NEW_VECTOR = 18
  NEW_LIST = 19
  DEFINE = 20
  CONS = 21
  INSERT = 22

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
      raise "No args can be nil"
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
      args = @args.zip(@debug).map { |a, d| "#{d || '?'} (#{a})" }.join("\t")
    else
      args = @args.join("\t")
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
      buffer << arg
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

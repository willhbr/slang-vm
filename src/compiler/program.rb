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
  def self.sizeof(type)
    case type
    when :local, :integer, :string, :atom, :global, :arg_count, :offset, :position
      1
    when :big_integer
      8
    else
      raise "Unknown type #{type}"
    end
  end

  def self.fits_in(type, value)
    case type
    when :integer
      0 <= value && value <= 255
    when :big_integer
      # TODO Check this
      true
    when :string, :atom, :global, :local
      0 <= value && value <= 255
    when :arg_count, :offset, :position
      0 <= value && value <= 255
    else
      raise "Unknown type #{type}"
    end
  end

  def self.write_bytes(buffer, type, value)
    case type
    when :local, :integer, :string, :atom, :global, :arg_count, :offset, :position
      format = 'C'
    when :big_integer
      format = 'Q>'
    else
      raise "Unknown type: #{type}"
    end
    [ value ].pack(format).split('').each do |byte|
      buffer << byte.ord
    end
  end

  INSTRUCTIONS = {
    LOAD_LOCAL: {
      id: 1,
      args: [:local]
    },
    LOAD_DEF: {
     id: 2,
     args: [:global]
    },
    STORE: {
     id: 3,
     args: [:local]
    },
    INVOKE: {
     id: 4,
     args: [:arg_count]
    },
    CONST_I: {
     id: 5,
     args: [:integer]
    },
    CONST_I_BIG: {
     id: 6,
     args: [:big_integer]
    },
    CONST_S: {
     id: 7,
     args: [:string]
    },
    CONST_A: {
     id: 8,
     args: [:atom]
    },
    CONST_TRUE: {
     id: 9,
     args: []
    },
    CONST_FALSE: {
     id: 10,
     args: []
    },
    CONST_NIL: {
     id: 11,
     args: []
    },
    JUMP: {
     id: 12,
     args: [:offset]
    },
    JUMP_BACK: {
     id: 13,
     args: [:offset]
    },
    AND: {
     id: 14,
     args: [:offset]
    },
    RETURN: {
     id: 15,
     args: []
    },
    CLOSURE: {
     id: 16,
     args: [:position, :arg_count],
     vararg: :local
    },
    PROTOCOL_CLOSURE: {
     id: 17,
     args: [:global]
    },
    NEW_MAP: {
     id: 18,
     args: [:arg_count]
    },
    NEW_VECTOR: {
     id: 19,
     args: [:arg_count]
    },
    NEW_LIST: {
     id: 20,
     args: [:arg_count]
    },
    DEFINE: {
     id: 21,
     args: [:global]
    },
    TYPE: {
     id: 22,
     args: [:global, :string, :arg_count],
     vararg: :local
    },
    INSTANCE: {
     id: 23,
     args: [:global, :arg_count]
    },
    IMPLEMENT: {
     id: 24,
     args: []
    },
    RAISE: {
     id: 25,
     args: []
    },
    TRY: {
     id: 26,
     args: [:offset]
    },
    END_TRY: {
     id: 27,
     args: []
    },
    DISCARD: {
     id: 28,
     args: []
    }
  }

  CODE_VALUES = Hash.new
  CODE_NAMES = Hash.new

  class << self
    Code::INSTRUCTIONS.each do |instruction, props|
      val = props[:id]
      Code::CODE_VALUES[instruction] = val
      Code::CODE_NAMES[val] = instruction
      define_method(instruction) do |args=[], debug=nil|
        args = args.is_a?(Array) ? args : [args]
        if debug != nil
          debug = debug.is_a?(Array) ? debug : [debug]
        end
        size = 1
        args.zip(props[:args]).map do |splat|
          arg, type = splat
          type ||= props[:vararg]
          raise "Too many arguments for #{instruction}: #{args}" unless type
          raise "Cannot have nil value: #{args}" unless arg
          size += sizeof(type)
          unless arg == :TEMP || fits_in(type, arg)
            raise "Value too big for #{type}: #{arg}"
          end
        end
        Code.new val, args, debug, size, props[:args], props[:vararg]
      end
    end
  end

  attr_reader :args
  attr_reader :debug
  def initialize(code, args, debug, size, arg_types, vararg_type)
    @code = code
    @args = args
    @debug = debug
    @size = size
    @arg_types = arg_types
    @vararg_type = vararg_type
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
    @args.zip(@arg_types).each do |splat|
      arg, type = splat
      type ||= @vararg_type
      self.class.write_bytes(buffer, type, arg)
    end
  end

  def size
    @size
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

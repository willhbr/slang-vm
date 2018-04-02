class Program
  attr_reader :position

  def initialize
    @strings = Hash.new # string to const id
    @string_code = 0
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

  def to_s
    @buffer.map(&:to_s).join("\n")
  end
end


class Code
  LOAD = 1
  STORE = 2
  DISPATCH = 3
  APPLY = 4
  CONST_I = 5
  CONST_S = 6
  JUMP = 7
  AND = 8
  OR = 9
  RETURN = 10
  # ===
  NEW_MAP = 11
  NEW_VECTOR = 12
  NEW_LIST = 13

  CONS = 14
  INSERT = 15

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
  def initialize(code, *args)
    @code = code
    @args = args
  end

  def to_s
    "#{self.class.stringify(@code)}\t#{@args.join("\t")}"
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

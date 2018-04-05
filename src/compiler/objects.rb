
class Vector < Array
  def inspect
    "[#{map(&:inspect).join ' '}]"
  end
end

class Array
  def inspect
    "(#{map(&:inspect).join ' '})"
  end
end

class Identifier
  attr_accessor :module
  attr_accessor :var
  attr_accessor :code
  attr_accessor :location

  KEYWORDS = [
    'let',
    'def',
    'fn',
    'true',
    'false',
    'nil'
  ]

  def initialize(value, location)
    sections = value.split('.')
    @module = sections[0..-2].join('.')
    @module = nil if @module == ''
    @var = sections[-1]
    @location = location
    @just_module = false
  end

  def to_module_only!
    @just_module = true
    if @module
      @var = "#{@module}.#{@var}"
      @module = nil
    end
  end

  def value
    if @module
      @module + '.' + @var
    else
      @var
    end
  end

  def local?
    @module == nil && !@just_module
  end

  def just_module?
    @just_module
  end

  def name_and_location
    line, col = @location
    "#{value} [#{line}:#{col}]"
  end


  def inspect
    if KEYWORDS.include? @value
      value
    else
      "#{value}_#{@code}"
    end
  end
end

class Atom
  attr_accessor :value
  attr_accessor :kw_arg

  def initialize(value, kw_arg=false)
    @value = value
    @kw_arg = kw_arg
  end

  def inspect
    @kw_arg ? "#{@value}:" : ":#{@value}"
  end
end

class Hash
  def inspect
    "{#{map { |k, v| "#{k.inspect} #{v.inspect}" }.join ' '}}"
  end
end

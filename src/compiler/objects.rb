
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
  attr_accessor :value
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
    @value = value
    @location = location
  end

  def name_and_location
    line, col = @location
    "#{@value} [#{line}:#{col}]"
  end


  def inspect
    if KEYWORDS.include? @value
      @value
    else
      "#{@value}_#{@code}"
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

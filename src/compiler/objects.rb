
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

  KEYWORDS = [
    'let',
    'def',
    'fn'
  ]

  def initialize(value)
    @value = value
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

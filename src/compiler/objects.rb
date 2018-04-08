class ClosureArgs
  attr_accessor :vars

  def initialize(vars)
    @vars = vars
  end

  def inspect
    '&[' + @vars.map(&:inspect).join(' ') + ']'
  end
end

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
  attr_accessor :code
  attr_accessor :location

  KEYWORDS = [
    'let',
    'if',
    'def',
    'fn',
    'true',
    'false',
    'nil',
    'spawn'
  ]

  def initialize(value, location)
    @value = value
    @parts = value.split('.')
    @location = location
    @is_module = @parts.length > 1
  end

  def whole
    @value
  end

  def module_part
    p = @parts[0..-2]
    if p.empty?
      nil
    else
      p.join '.'
    end
  end

  def local?
    !@is_module
  end

  def global?
    @is_module
  end

  # For one-word aliases and whatnot
  def make_global!
    @is_module = true
  end

  def no_module?
    @parts.size == 1
  end

  def var_part
    @parts[-1]
  end

  def parts
    [module_part, var_part]
  end

  def add_module!(mod)
    @parts.insert(0, mod)
    @value = "#{mod}.#{@value}"
  end

  def name_and_location
    line, col = @location
    "#{@value} [#{line}:#{col}]"
  end


  def inspect
    if KEYWORDS.include?(@value) || @code.nil?
      @value
    else
      "#{@value}_#{@code}"
    end
  end
end

class Atom
  attr_accessor :code
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

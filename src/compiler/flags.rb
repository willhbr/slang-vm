class Flags
  @@flags = Hash.new
  @@args = Array.new

  def self.[](key)
    @@flags[key]
  end

  def self.[]=(key, value)
    @@flags[key] = value
  end

  def self.<<(arg)
    @@args << arg
  end

  def self.args
    @@args
  end

  def self.to_s
    @@flags.map { |k, v| "--#{k}=#{v}" }.join("\n")
  end
end

ARGV.each do |arg|
  if arg.start_with? '--'
    key, eq, value = arg[2..-1].partition '='
    value = true if eq == ''
    if value == 'true'
      value = true
    elsif value == 'false'
      value = false
    end
    Flags[key.to_sym] = value
    if key.start_with? 'no'
      Flags[key[2..-1].to_sym] = false
    end
  else
    Flags << arg
  end
end

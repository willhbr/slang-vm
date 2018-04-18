require_relative './builtins'

class Def
  def generate_def(file)
    file.puts "// #{self.to_s}: #{@code}"
    case @type
    when :type
      file.puts "#{@name}Type,"
    when :module
      file.puts "Module{Name: \"#{@name}\"},"
    when :func, :impl
      file.puts "GoClosure{Function: #{@module}__#{to_go_name(@name.to_s)}},"
    when :protocol
      file.puts "ProtocolClosure{ID: #{@code}},"
    else
      nil
    end
  end

  def generate_type_with_impls(file)
    file.puts "// #{self.to_s}: #{@code}"
    file.puts "var #{@name}Type = &Type{Name: \"#{@name}\","
    file.puts "ProtocolMethods: map[int]Closure{"
    @children.select { |c| c.type == :impl }.each do |d|
      file.puts "#{d.implements.code}: GoClosure{Function: #{@name}__#{to_go_name(d.name.to_s)}},"
    end
    file.puts '}}'
  end
end

def to_go_name(str)
  {
    '->' => '_rArr_',
    '<-' => '_lArr_',
    '<' => 'lessThan',
    '-' => 'minus',
    '*' => 'times',
    '=' => '_eq_',
    '+' => '_plus_',
    '/' => '_div_'
  }.reduce str do |str, replace|
    str.gsub(*replace)
  end
end

if __FILE__==$0
  File.open(ARGV[0], 'w') do |file|
    file.puts "package types"
    file.puts 'var Defs = []Value {'
    Def::BUILTIN_DEFS.sort_by { |d| d.code }.each do |d|
      d.generate_def(file)
    end
    file.puts '}'
  end

  File.open(ARGV[1], 'w') do |file|
    file.puts "package types"
    file.puts 'import "math/big"'
    Def::BUILTIN_DEFS.select { |d| d.type == :type }.each do |d|
      d.generate_type_with_impls(file)
    end

    file.puts """
    func GetType(object Value) *Type {
      switch object.(type) {
    """
    Def::BUILTIN_DEFS.select { |d| d.type == :type }.each do |d|
      if d.gotype
        file.puts "case #{d.gotype}:"
      else
        file.puts "case #{d.name}:"
      end
      file.puts "return #{d.name}Type"
    end
    file.puts """
    case Instance:
      return object.(Instance).Type
    default:
      println(object)
      panic(\"Can't find type of object\")
    }}
    """
  end
end

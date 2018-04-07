require_relative './objects'

class Builtins
  MODULES = {
    IO: [
      :puts,
      :gets
    ],
    Kernel: [
      :type,
      :<,
      :-,
      :*
    ],
    Channel: [
      :new,
      :send,
      :receive
    ]
  }
end

class Defs
  @@defs = Hash.new { |h, k| h[k] = Hash.new }

  @@next_code = 0
  def self.defs
    @@defs
  end

  def self.def_count
    @@next_code
  end

  def self.alias(current, from, to)
    @@defs[current] = @@defs[]
  end

  def self.get?(identifier, mod=nil)
    if identifier.module
      if r = @@defs[identifier.module][identifier.var]
        return r
      end
    else
      if r = @@defs[mod][identifier.var]
        return r
      end
    end
    raise "Unknown module for #{identifier.value} (have you declared a module?)" unless mod
    nil
  end

  def self.set_and_return(identifier, mod=nil)
    identifier.module ||= mod
    raise "Unknown module for #{identifier.value} (have you declared a module?)" unless identifier.module
    if iden = @@defs[identifier.module][identifier.var]
      identifier.code = iden.code
    else
      identifier.code = @@next_code
      @@defs[identifier.module][identifier.var] = identifier
      @@next_code += 1
    end
    identifier.code
  end

  def self.define_module(identifier)
    identifier.to_module_only!
    if iden = @@defs[identifier.var][:__MODULE__]
      identifier.code = iden.code
    else
      identifier.code = @@next_code
      @@defs[identifier.var][:__MODULE__] = identifier
      @@next_code += 1
    end
    identifier.code
  end
end

Builtins::MODULES.sort.each do |name, methods|
  Defs.define_module(Identifier.new(name.to_s, [nil, nil]))
  methods.sort.each do |method|
    Defs.set_and_return(Identifier.new(method.to_s, [nil, nil]), name.to_s)
  end
end

def to_go_name(str)
  {
    '->' => 'arrow',
    '<-' => 'reverseArrow',
    '<' => 'lessThan',
    '-' => 'minus',
    '*' => 'times'
  }.reduce str do |str, replace|
    str.gsub(*replace)
  end
end

if __FILE__==$0
  File.open(ARGV[0], 'w') do |file|
    file.puts "package funcs"
    file.puts 'import "../ds"'
    file.puts 'var Defs = []ds.Value {'

    Defs.defs.each do |mod, defs|
      module_identifier = defs[:__MODULE__]
      file.puts "// #{mod}: #{module_identifier.code}"
      file.puts "ds.Module{Name: \"#{module_identifier.value}\"},"
      defs.each do |name, iden|
        next if name == :__MODULE__
        file.puts "// #{iden.value}: #{iden.code}"
        file.puts "GoClosure{Function: #{iden.module}__#{to_go_name(iden.var)}},"
      end
    end
    file.puts '}'
  end
end

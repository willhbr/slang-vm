require_relative './objects'

class Defs
  @@defs = Hash.new
  @@next_code = 0
  def self.defs
    @@defs
  end

  def self.def_count
    @@next_code
  end

  def self.get?(identifier, mod=nil)
    if mod && (in_this = @@defs[mod + '.' + identifier.value])
      return in_this
    end
    if iden = @@defs[identifier.value]
      return iden
    end
    raise "Unknown module for #{identifier.value} (have you declared a module?)" unless mod
    nil
  end

  def self.set_and_return(identifier, mod=nil)
    identifier.module ||= mod
    raise "Unknown module for #{identifier.value} (have you declared a module?)" unless identifier.module
    if iden = @@defs[identifier.value]
      identifier.code = iden.code
    else
      identifier.code = @@next_code
      @@defs[identifier.value] = identifier
      @@next_code += 1
    end
    identifier.code
  end

  def self.define_module(identifier)
    identifier.to_module_only!
    if iden = @@defs[identifier.value]
      identifier.code = iden.code
    else
      identifier.code = @@next_code
      @@defs[identifier.value] = identifier
      @@next_code += 1
    end
    identifier.code
  end
end

class Builtins
  MODULES = {
    IO: [
      :puts,
    ],
    Kernel: [
      :type
    ]
  }
end

Builtins::MODULES.sort.each do |name, methods|
  Defs.define_module(Identifier.new(name.to_s, [nil, nil]))
  methods.sort.each do |method|
    Defs.set_and_return(Identifier.new(method.to_s, [nil, nil]), name.to_s)
  end
end

if __FILE__==$0
  File.open(ARGV[0], 'w') do |file|
    file.puts "package funcs"
    file.puts 'import "../ds"'

    file.puts 'var Defs = []ds.Value {'

    Defs.defs.each do |name, iden|
      file.puts "// #{iden.value}: #{iden.code}"
      if iden.just_module? # It's a module literal
        file.puts "ds.Module{Name: \"#{iden.value}\"},"
      else
        file.puts "GoClosure{Function: #{iden.module}__#{iden.var.gsub('-', '_')}},"
      end
    end
    file.puts '}'
  end
end

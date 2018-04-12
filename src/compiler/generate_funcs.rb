require_relative './defs'

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
      file.puts "ds.Module{Name: \"#{module_identifier.whole}\"},"
      defs.each do |name, iden|
        next if name.is_a?(Symbol)
        if (proto = Builtins::PROTOCOL_METHODS[mod.to_sym]) && (proto.include? name.to_sym)
          if iden.location.nil? # It's an internal method, not overriden
            file.puts "// #{iden.whole}: #{iden.code} (proto method)"
            file.puts "ProtocolClosure{ID: #{iden.code}},"
            next
          end
        end

        file.puts "// #{iden.whole}: #{iden.code}"
        file.puts "GoClosure{Function: #{iden.module_part}__#{to_go_name(iden.var_part)}},"
      end
    end
    file.puts '}'
  end
end

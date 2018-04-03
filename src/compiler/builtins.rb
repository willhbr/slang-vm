class Builtins
  MODULES = {
    IO: [
      :puts
    ]
  }

  class << self
    next_module = 0
    with_codes = Hash.new
    Builtins::MODULES.sort.each do |mod, methods|
      method_codes = Hash.new
      methods.sort.each_with_index do |method, index|
        method_codes[method] = index
        method_codes[index] = method
      end
      with_codes[mod] = {
        id: next_module,
        name: mod,
        methods: method_codes,
      }
      with_codes[next_module] = with_codes[mod]
      next_module += 1
    end

    MODULE_CODES = with_codes

    def [](mod, method)
      mod_methods = MODULE_CODES[mod]
      return unless mod_methods
      mod_id = mod_methods[:id]
      method_id = mod_methods[:methods][method]
      return unless method_id
      [mod_id, method_id]
    end

    def module_code(mod)
      mod_methods = MODULE_CODES[mod]
      mod_methods[:id]
    end

    def module_name(mod)
      mod_methods = MODULE_CODES[mod]
      mod_methods[:name]
    end

    def modules
      MODULE_CODES
    end
  end
end

if __FILE__==$0
  Builtins.modules
  File.open(ARGV[0], 'w') do |file|
    file.puts "package funcs"
    file.puts 'import "../ds"'

    file.puts 'var Modules = [][]func(ds.Value) ds.Value{'

    Builtins.modules.each do |name, info|
      next unless name.is_a? Symbol
      file.puts "{"
      info[:methods].each do |method, id|
        next unless method.is_a? Symbol
        file.puts "// #{name}__#{method}: #{info[:id]}"
        file.puts "#{name}__#{method.to_s.gsub('-', '_')},"
      end
      file.puts '},'
    end
    file.puts '}'
  end
end

require 'time'

require_relative './flags'
require_relative './local_resolver'
require_relative './def_resolver'
require_relative './objects'
require_relative './program'
require_relative './parser'
require_relative './code_generator'
require_relative './macro_expander'
require_relative './dag_order'

def read_file(filename)
  scanner = Scanner.new filename, File.read(filename)
  tokens = scanner.read
  ast = Parser.new(tokens).program
  if Flags[:show_ast]
    puts "===#{filename}==="
    puts ast.inspect
  end
  ast
end

bad_macro_expander = MacroExpander.new
local = LocalResolver.new
global = DefResolver.new
code_gen = CodeGenerator.new


start = Time.now
Flags.args[1..-1].each do |path|
  tree = read_file(path)
  tree = tree.map do |node|
    bad_macro_expander.process_top_level(node)
  end

  if Flags[:show_expanded]
    puts "===#{path}==="
    puts tree.inspect
  end

  tree.each do |node|
    local.process_top_level(node)
  end

  tree.each do |node|
    global.process_top_level(node)
  end

  if Flags[:show_resolved]
    puts "===#{path}==="
    puts tree.inspect
  end
end

ordered, module_order = DAGOrder.order(global.modules, global.module_useage)

if Flags[:show_modules]
  puts "===Modules==="
  puts global.modules.keys.join(', ')
end
if Flags[:show_order]
  puts "===Order==="
  puts module_order.join(', ')
end
if Flags[:show_ordered_ast]
  puts "===Ordered AST==="
  puts ordered.inspect
end


ordered.map do |node|
  code_gen.process_top_level(node)
end

if Flags[:show_instructions]
  puts "===Instructions==="
  code_gen.program.print
end

fin_time = Time.now

if Flags[:show_time]
  puts "===Duration==="
  puts "#{(fin_time - start) * 1000.0}ms"
end

File.open(Flags.args[0], 'wb') do |output|
  code_gen.program.write_to(output)
end

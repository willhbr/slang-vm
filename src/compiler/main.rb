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
  Parser.new(tokens).program
end

bad_macro_expander = MacroExpander.new
local = LocalResolver.new
global = DefResolver.new
code_gen = CodeGenerator.new


ARGV[1..-1].each do |path|
  tree = read_file(path)
  tree = tree.map do |node|
    bad_macro_expander.process_top_level(node)
  end

  tree.each do |node|
    local.process_top_level(node)
  end

  tree.each do |node|
    global.process_top_level(node)
  end
end

ordered = DAGOrder.order(global.modules, global.module_useage)

ordered.map do |node|
  code_gen.process_top_level(node)
end

File.open(ARGV[0], 'wb') do |output|
  code_gen.program.write_to(output)
end

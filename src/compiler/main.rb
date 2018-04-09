require_relative './resolve'
require_relative './objects'
require_relative './program'
require_relative './parser'
require_relative './code_generator'
require_relative './macro_expander'

filename = ARGV[0]
scanner = Scanner.new filename, File.read(filename)
tokens = scanner.read
tree = Parser.new(tokens).program

bad_macro_expander = MacroExpander.new

tree = tree.map do |node|
  bad_macro_expander.process_top_level(node)
end

res = Resolver.new

tree.each do |node|
  res.process_top_level(node)
end

p tree

cg = CodeGenerator.new
tree.map do |node|
  cg.process_top_level(node)
end

cg.program.print

File.open(ARGV[1], 'wb') do |output|
  cg.program.write_to(output)
end

require_relative './resolve'
require_relative './objects'
require_relative './program'
require_relative './parser'
require_relative './code_generator'

filename = ARGV[0]
scanner = Scanner.new filename, File.read(filename)
tokens = scanner.read
tree = Parser.new(tokens).program

res = Resolver.new()

tree.map do |node|
  res.resolve_top_level(node)
end

cg = CodeGenerator.new
tree.map do |node|
  cg.generate(node)
end

p tree
puts cg.program

File.open(ARGV[1], 'wb') do |output|
  cg.program.write_to(output)
end

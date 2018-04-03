require_relative './objects'
require_relative './program'

class CodeGenerator
  attr_reader :program
  def initialize
    @program = Program.new
  end

  def generate(ast)
    case ast
    when Vector
      push Code.NEW_VECTOR
      ast.each do |node|
        generate(node)
        push Code.CONS
      end
    when Array
      generate_call(ast)
    when Hash
      push Code.NEW_MAP
      ast.each do |key, value|
        generate(key)
        generate(value)
        push Code.INSERT
      end
    when Identifier
      case ast.value
      when 'true'
        push Code.CONST_TRUE
      when 'true'
        push Code.CONST_FALSE
      when 'nil'
        push Code.CONST_NIL
      else
        push Code.LOAD(ast.code)
      end
    when String
      code = @program.add_string(ast)      
      push Code.CONST_S(code, ast)
    when Integer
      push Code.CONST_I(ast)
    end
  end

  def generate_call(ast)
    first = ast.first
    unless first
      push Code.NEW_LIST
    end
    if first.is_a?(Identifier) && first.value == 'let'
      binds = ast[1]
      binds.each_slice(2) do |slice|
        name, expr = slice
        raise "Not an identifier: #{name}" unless name.is_a? Identifier
        generate(expr)
        push Code.STORE(name.code, name.name_and_location)
      end
      ast[2..-1].each do |node|
        generate(node)
      end
    elsif first.is_a?(Identifier) && first.value == 'def'
      nil
    elsif first.is_a?(Identifier) && first.value == 'if'
      _, cond, then_block, else_block = ast
      generate(cond)
      jump = Code.AND(-1) # Don't know where to jump to
      push jump
      start = @program.position
      generate(then_block)
      if else_block
        else_jump = Code.JUMP(-1)
        push else_jump
        else_start = @program.position
        jump.args[0] = @program.position - start
        generate(else_block)
        # jump to end of else block
        else_jump.args[0] = @program.position - else_start
      else
        # jump to end of if
        jump.args[0] = @program.position - start
      end
    elsif first.is_a?(Identifier)
      ast[1..-1].each do |arg|
        generate(arg)
      end
      push Code.DISPATCH(first.code, first.name_and_location)
    else
      ast.each do |arg|
        generate(arg)
      end
      push Code.DISPATCH(-1, 'apply')
    end
  end

  def push(*args)
    @program.<<(*args)
  end
end

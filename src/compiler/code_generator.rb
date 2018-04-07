require_relative './objects'
require_relative './program'

class CodeGenerator
  attr_reader :program
  def initialize
    @program = Program.new
  end

  def generate_top_level(ast)
    @func_recur_points = []
    generate(ast)
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
      when 'false'
        push Code.CONST_FALSE
      when 'nil'
        push Code.CONST_NIL
      else
        if ast.local?
          push Code.LOAD_LOCAL(ast.code, ast.value)
        else
          push Code.LOAD_DEF(ast.code, ast.value)
        end
      end
    when String
      code = @program.add_string(ast)      
      push Code.CONST_S(code, ast)
    when Integer
      push Code.CONST_I(ast)
    when Atom
      code = @program.add_string(ast.value)
      push Code.CONST_A(code, ast.value)
    end
  end

  def generate_call(ast)
    first = ast.first
    unless first
      push Code.NEW_LIST
    end
    if first.is_a? Identifier
      case first.value
      when 'let'
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
      when 'module'
        name = ast[1]
        str = @program.add_string(name.value)
        # TODO make this do a module
        push Code.CONST_S(str, name.value)
        push Code.DEFINE(name.code, name.name_and_location)
      when 'def'
        name = ast[1]
        expr = ast[2]
        raise "def name must be identifier" unless name.is_a? Identifier
        raise 'def accepts 2 args' if ast.size > 3
        generate(expr)
        push Code.DEFINE(name.code, name.name_and_location)
      when 'spawn'
        spawn = Code.SPAWN(-1)
        push spawn
        idx = @program.position
        generate(ast[1])
        spawn.args[0] = @program.position - idx
      when 'recur'
        ast[1..-1].each do |expr|
          generate(expr)
        end
        # Extra two for the two instructions pushed
        push Code.JUMP_BACK(@program.position - @func_recur_points[-1] + 2)
      when 'fn'
        jump = Code.JUMP(-1)
        push jump
        pos = @program.position
        @func_recur_points << pos
        args = ast[1]
        captured = args.pop
        args.reverse.each do |arg|
          push Code.STORE(arg.code, arg.value)
        end

        body = ast[2..-1]
        body.each do |node|
          generate(node)
        end
        @func_recur_points.pop
        push Code.RETURN
        jump.args[0] = @program.position - pos
        closure_args = [pos, captured.vars.length]
        debug = [nil, nil]
        captured.vars.each do |iden|
          closure_args << iden.code
          debug << iden.value
        end
        push Code.CLOSURE(closure_args, debug)
      when 'if'
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
      when 'do'
        ast[1..-1].each do |arg|
          generate(arg)
        end
      else
        ast[1..-1].each do |arg|
          generate(arg)
        end
        generate(ast[0])
        # Minus one for function name
        arg_count = ast.size - 1
        push Code.INVOKE(arg_count, ast[0].is_a?(Identifier) ? ast[0].value : nil)
      end
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

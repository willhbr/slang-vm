require_relative './objects'
require_relative './program'
require_relative './ast_processor'

class CodeGenerator
  include ASTProcessor

  attr_reader :program
  def initialize
    @program = Program.new
  end

  def reset
    @func_recur_points = []
  end

  def process_vector(ast, top_level)
    ast.each do |node|
      process(node)
    end
    push Code.NEW_VECTOR(ast.size)
  end

  def process_array(ast, top_level)
    generate_call(ast)
  end

  def process_hash(ast, top_level)
    ast.each do |key, value|
      process(key)
      process(value)
    end
    push Code.NEW_MAP(ast.size)
  end

  def process_identifier(ast, top_level)
    case ast.whole
    when 'true'
      push Code.CONST_TRUE
    when 'false'
      push Code.CONST_FALSE
    when 'nil'
      push Code.CONST_NIL
    else
      if ast.local?
        push Code.LOAD_LOCAL(ast.code, ast.whole)
      else
        push Code.LOAD_DEF(ast.code, ast.whole)
      end
    end
  end

  def process_string(ast, top_level)
    code = @program.add_string(ast)      
    push Code.CONST_S(code, ast)
  end

  def process_atom(ast, top_level)
    # TODO Use other set of IDs for atoms
    code = new_atom(ast)
    push Code.CONST_A(code, ast)
  end

  def new_atom(ast)
    @program.add_string(ast.value)
  end

  def process_integer(ast, top_level)
    if ast > 255 || ast < 0
      # Encode as big endian, 64-bit (8 bytes) signed
      # TODO do this in a non-shitty way
      push Code.CONST_I_BIG([ ast ].pack('Q>').split('').map(&:ord), ast.to_s)
    else
      push Code.CONST_I(ast)
    end
  end

  def process_atom(ast, top_level)
    code = @program.add_string(ast.value)
    push Code.CONST_A(code, ast.value)
  end

  def process_array(ast, top_level)
    first = ast.first
    unless first
      push Code.NEW_LIST
    end
    if first.is_a? Identifier
      case first.whole
      when 'let'
        binds = ast[1]
        binds.each_slice(2) do |slice|
          name, expr = slice
          raise "Not an identifier: #{name}" unless name.is_a? Identifier
          process(expr)
          push Code.STORE(name.code, name.name_and_location)
        end
        ast[2..-1].each do |node|
          process(node)
        end
      when 'alias', 'import'
        return
      when 'module'
        name = ast[1]
        str = @program.add_string(name.whole)
        # TODO make this do a module
        push Code.CONST_S(str, name.whole)
        push Code.DEFINE(name.code, name.name_and_location)
      when 'def'
        name = ast[1]
        expr = ast[2]
        raise "def name must be identifier" unless name.is_a? Identifier
        raise 'def accepts 2 args' if ast.size > 3
        process(expr)
        push Code.DEFINE(name.code, name.name_and_location)
      when 'recur'
        ast[1..-1].each do |expr|
          process(expr)
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
          push Code.STORE(arg.code, arg.whole)
        end

        body = ast[2..-1]
        body.each do |node|
          process(node)
        end
        @func_recur_points.pop
        push Code.RETURN
        jump.args[0] = @program.position - pos
        closure_args = [pos, captured.vars.length]
        debug = [nil, nil]
        captured.vars.each do |iden|
          closure_args << iden.code
          debug << iden.whole
        end
        push Code.CLOSURE(closure_args, debug)
      when 'if'
        _, cond, then_block, else_block = ast
        process(cond)
        jump = Code.AND(-1) # Don't know where to jump to
        push jump
        start = @program.position
        process(then_block)
        if else_block
          else_jump = Code.JUMP(-1)
          push else_jump
          else_start = @program.position
          jump.args[0] = @program.position - start
          process(else_block)
          # jump to end of else block
          else_jump.args[0] = @program.position - else_start
        else
          # jump to end of if
          jump.args[0] = @program.position - start
        end
      when 'deftype'
        name = ast[1]
        attrs = ast[2..-1]
        str = @program.add_string(name.whole)
        codes = [name.code, str]
        debug = [name.whole, '-', 'size']
        codes.push attrs.size
        attrs.each do |attr|
          code = new_atom(attr)
          debug.push attr.value
          codes.push code
        end
        push Code.TYPE(codes, debug)
      when 'new-instance'
        name = ast[1]
        attrs = ast[2..-1]
        attrs.each do |attr|
          process(attr)
        end
        push Code.INSTANCE([name.code, attrs.size], [name.whole, 'size'])
      when 'do'
        ast[1..-1].each do |arg|
          process(arg)
        end
      else
        ast[1..-1].each do |arg|
          process(arg)
        end
        process(ast[0])
        # Minus one for function name
        arg_count = ast.size - 1
        push Code.INVOKE(arg_count, ast[0].is_a?(Identifier) ? ast[0].whole : nil)
      end
    else
      puts "Woops: #{ast.inspect}"
      ast.each do |arg|
        process(arg)
      end
      push Code.INVOKE(-1, 'apply')
    end
  end

  def push(*args)
    @program.<<(*args)
  end
end

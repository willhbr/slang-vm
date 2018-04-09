require_relative './objects'
require_relative './builtins'
require_relative './ast_processor'

class ResolverState
  attr_accessor :scopes
  attr_accessor :in_func

  def initialize
    @scopes = [ Hash.new ]
    @in_func = 0
  end

  def rescope(&block)
    @scopes.push Hash.new
    result = block.()
    @scopes.pop
    result
  end
end

class LocalResolver
  include ASTProcessor
  def initialize
    @next_id = 0
    @on_local_bind = []
    @state = ResolverState.new
  end

  def reset
    @next_id = 0
  end

  def process_other(ast, top_level)
    passthrough_nested(ast)
  end

  def process_identifier(ast, top_level)
    return if Identifier::KEYWORDS.include? ast.whole
    return if bind_to_local(ast)
  end

  def process_fn(ast)
    @state.in_func += 1
    @state.rescope do
      args = ast[1]
      outside_code = @next_id
      raise "Args must be vector, not #{args.class} #{ast[0].location}" unless args.is_a? Vector
      args.each do |arg|
        raise "Arg must be identifier: #{arg}" unless arg.is_a? Identifier
        bind_new!(arg)
      end
      captured = []
      # This is a hackety hack
      @on_local_bind << lambda do |iden|
        if iden.code < outside_code
          captured << iden
        end
      end
      args.push(ClosureArgs.new captured)
      ast[2..-1].each do |expr|
        process(expr)
      end
      @on_local_bind.pop
    end
    @state.in_func -= 1
  end

  def process_array(ast, top_level)
    first = ast.first
    return unless first
    if first.is_a? Identifier
      case first.whole
      when 'let'
        binds = ast[1]
        @state.rescope do
          binds.each_slice(2) do |slice|
            name, expr = slice
            process(expr)
            raise "Not an identifier: #{name}" unless name.is_a? Identifier
            bind_new!(name)
          end
          ast[2..-1].each do |node|
            process(node)
          end
        end
      # These are skipped to be resolved later
      when 'recur'
        if @state.in_func != 0
          ast[1..-1].each { |node| process(node, top_level) }
        else
          raise "Can't recur outside function"
        end
      when 'fn'
        process_fn(ast)
      else
        if Identifier::KEYWORDS.include? first.whole
          ast[1..-1].each { |node| process(node) }
        else
          ast.each { |node| process(node) }
        end
      end
    else
      ast.each { |node| process(node) }
    end
  end

  def bind_to_local(identifier)
    @state.scopes.reverse.each do |scope|
      if bound = scope[identifier.whole]
        identifier.code = bound.code
        @on_local_bind.each do |func|
          func.(identifier)
        end
        return true
      end
    end
    return false
  end

  def bind_new!(identifier)
    identifier.code = @next_id
    @next_id += 1
    @state.scopes[-1][identifier.whole] = identifier
  end
end

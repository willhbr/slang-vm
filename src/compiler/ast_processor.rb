require_relative './flags'

module ASTProcessor
  def process_top_level(ast)
    reset()
    res = process(ast, true)
    done_statement(ast)
    res
  end

  def reset
  end

  def done_statement(ast)
  end

  private

  def process(ast, top_level=false)
    if Flags[:show_process]
      puts "#{('[ ' + self.class.to_s + ' ]').ljust(20, ' ')} #{ast.inspect}"
    end
    name = ast.class.to_s.downcase
    if name == 'fixnum'
      # Fix for ruby 2.3
      name = 'integer'
    end
    if self.respond_to? :"process_#{name}"
      self.send(:"process_#{name}", ast, top_level)
    elsif self.respond_to? :process_other
      self.send(:process_other, ast, top_level)
    else
      raise """
      No process implemented for #{name}
      Compiling: #{ast.inspect}
      """
    end
  end

  def passthrough_nested(ast)
    case ast
    when Vector
      Vector.new(ast.map do |node|
        process(node)
      end)
    when Array
      ast.map do |node|
        process(node)
      end
    when Hash
      result = Hash.new
      ast.each do |key, value|
        result[process(key)] = process(value)
      end
      result
    else
      ast
    end
  end
end

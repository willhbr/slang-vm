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
    name = ast.class.to_s.downcase
    if self.respond_to? :"process_#{name}"
      self.send(:"process_#{name}", ast, top_level)
    else
      self.send(:process_other, ast, top_level)
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

require_relative './objects'
require_relative './builtins'
require_relative './ast_processor'

class DefResolver
  include ASTProcessor

  def initialize
    @current_module = nil
  end

  def process_other(ast, top_level)
    passthrough_nested(ast)
  end

  def process_identifier(ast, top_level)
    return if Identifier::KEYWORDS.include? ast.whole
    return unless ast.code.nil?
    raise "No module defined!" unless @current_module
    return if Defs.get_module?(@current_module, ast)
    return if Defs.get_module_def?(@current_module, ast)
    raise "Undefined var: #{ast.whole} #{ast.location}"
  end

  def process_array(ast, top_level)
    first = ast.first
    return unless first
    if first.is_a? Identifier
      case first.whole
      when 'alias'
        raise 'cannot alias outside module' unless @current_module
        mod = ast[1]
        raise 'Can only alias modules' unless mod.is_a? Identifier
        check_unbound! mod
        process(mod)
        alias_to = ast[2]
        raise 'Can only alias modules' unless alias_to.is_a?(Identifier)
        check_unbound! alias_to
        Defs.alias(@current_module, mod, alias_to)
      when 'import'
        raise 'cannot import outside module' unless @current_module
        mod = ast[1]
        raise 'Can only alias modules' unless mod.is_a? Identifier
        check_unbound! mod
        process(mod)
        Defs.import(@current_module, mod)
      when 'module'
        raise "Can only define module at top-level" unless top_level
        _, name = ast
        check_unbound! name
        Defs.define_module(name)
        @current_module = name
      when 'def'
        raise "Can only def at top-level" unless top_level
        first, name, value = ast
        check_unbound! name
        process(value)
        raise 'Cannot define outside of module' unless @current_module
        raise "Cannot define in other module: #{name.mod}" unless name.no_module?
        if Defs.get_module_def?(@current_module, name)
          raise "Already defined #{name.whole} on #{iden.location}"
        end
        Defs.def_def(@current_module, name)
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

  def check_unbound!(iden)
    raise "#{iden.whole} is already bound as a local variable! #{iden.location}" unless iden.code.nil?
  end
end

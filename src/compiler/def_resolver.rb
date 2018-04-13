require_relative './objects'
require_relative './defs'
require_relative './ast_processor'

class DefResolver
  include ASTProcessor

  attr_reader :modules
  attr_reader :module_useage

  def initialize
    @current_module = nil
    @modules = Hash.new
    @current_list = nil
    @module_useage = Hash.new
  end

  def done_statement(ast)
    @current_list << ast
  end
 
  def process_other(ast, top_level)
    passthrough_nested(ast)
  end

  def process_identifier(ast, top_level)
    return if Identifier::KEYWORDS.include? ast.whole
    return unless ast.code.nil?
    raise "No module defined!" unless @current_module
    if Defs.get_module?(@current_module, ast)
      use! ast
      return
    end
    if splat = Defs.get_module_def?(@current_module, ast)
      d, mod = splat
      use! mod
      return
    end
    puts @current_module.inspect
    puts "Current defs: #{Defs.lookup_module_defs(@current_module, @current_module.whole).inspect}"
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
        define! name
        @current_module = name
      when 'def'
        raise "Can only def at top-level" unless top_level
        first, name, value = ast
        check_unbound! name
        process(value)
        raise 'Cannot define outside of module' unless @current_module
        if name.no_module?
          mod = @current_module
        else
          mod_iden = Identifier.from(name, name.module_part)
          mod = Defs.get_module? @current_module, mod_iden
          unless mod
            mod = Defs.define_module mod_iden
          end
        end
        if Defs.get_module_def?(mod, name)
          p mod, name
          raise "Already defined #{name.whole} on #{name.location}"
        end
        Defs.def_def(mod, name)
      when 'new-type'
        raise "Can only def at top-level" unless top_level
        name = ast[1]
        raise "Name must be identifier" unless name.is_a? Identifier
        attributes = ast[2..-1]
        raise "attributes must be atoms, got: #{attributes.map(&:class)}" unless attributes.all? { |a| a.is_a? Atom }        
        Defs.define_module(name) unless name == @current_module
        use! name
      when 'defprotocol'
        name = ast[1]
        raise "Protocol method name must be identifier" unless name.is_a? Identifier
        Defs.def_def @current_module, name
      when 'do'
        ast[1..-1].each do |node|
          process(node, top_level)
        end
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

  def define!(mod)
    @current_list = (@modules[mod.whole] ||= [])
    @current_useage = (@module_useage[mod.whole] ||= [])
  end

  def use!(other_mod)
    return if other_mod.code == @current_module.code
    @current_useage << other_mod.whole unless @current_useage.include? other_mod.whole
  end

  def check_unbound!(iden)
    raise "#{iden.whole} is already bound as a local variable! #{iden.location}" unless iden.code.nil?
  end
end

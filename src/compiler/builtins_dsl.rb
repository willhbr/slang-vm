require_relative './objects'
require_relative './defs'

class Def
  BUILTIN_DEFS = []
  @@protocols = Hash.new

  attr_accessor :module
  attr_accessor :implements
  attr_accessor :type
  attr_accessor :code
  attr_accessor :name

  def initialize(name)
    @name = name
    @module = nil
    @type = nil
    @code = nil
    @implements = nil
    @children = []
    BUILTIN_DEFS << self
  end

  def defn(name)
    raise "Can't define #{name} inside of #{@type}: #{@name}" unless @type == :module || @type == :type
    d = Def.new name
    d.type = :func
    @children << d
    d.module = self
  end

  def defimpl(name, options={})
    raise "Can't define #{name} inside of #{@type}: #{@name}" unless @type == :module || @type == :type
    d = Def.new name
    d.module = self
    d.type = :impl
    mod, method = options[:of]
    methods = @@protocols[mod]
    raise "No module #{mod}" unless methods
    iden = methods[method]
    raise "No protocol method #{method}" unless iden
    d.implements = iden
    @children << d
  end

  def defprotocol(name)
    raise "Can't define #{name} inside of #{@type}: #{@name}" unless @type == :module
    d = Def.new name
    d.module = self
    d.type = :protocol
    protos = (@@protocols[@name] ||= Hash.new)
    protos[name] = d
    @children << d
  end

  def self.assign_codes
    BUILTIN_DEFS.each do |d|
      if d.type == :module || d.type == :type
        iden = Defs.define_module d.to_identifier
      else
        iden = Defs.def_def d.module.to_identifier, d.to_identifier
      end
      d.code = iden.code
    end
  end

  def to_s
    if @module
      "#{@module}.#{@name}"
    else
      @name.to_s
    end
  end

  def to_identifier
    Identifier.new @name.to_s, nil
  end
end

def defmodule(name, &block)
  mod = Def.new name.to_sym
  mod.type = :module
  mod.instance_exec(&block)
end

def deftype(name, &block)
  mod = Def.new name.to_sym
  mod.type = :type
  mod.instance_exec(&block)
end



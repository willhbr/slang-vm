class Defs
  @@defs = Hash.new

  @@next_code = 0
  def self.defs
    @@defs
  end

  def self.def_count
    @@next_code
  end

  def self.alias(current, from, to)
    defs = @@defs[current.whole] ||= blank_def(current)
    from_mod = get_module?(current, from)
    from.code = from_mod.code
    to.code = from_mod.code
    from_defs = @@defs[from_mod.whole] ||= blank_def(from)
    defs[:__ALIAS__][to.whole] = from_defs
  end

  def self.import(current, imported)
    defs = @@defs[current.whole]
    to_import_mod = get_module?(current, imported)
    import_defs = @@defs[to_import_mod.whole]
    defs[:__IMPORTS__] << import_defs
  end

  def self.get_module?(current, identifier)
    defs = lookup_module_defs(current, identifier.whole)
    return unless defs
    mod = defs[:__MODULE__]
    identifier.code = mod.code
    identifier.make_global!
    mod
  end

  def self.lookup_module_defs(current, name)
    defs = @@defs[name]
    return defs if defs
    defs = @@defs[current.whole]
    aliased_mod = defs[:__ALIAS__][name]
    return aliased_mod if aliased_mod
  end

  def self.get_module_def?(current, identifier)
    mod, var = identifier.parts
    mod = current.whole if mod.nil?
    defs = lookup_module_defs current, mod
    return unless defs
    if v = defs[var]
      identifier.code = v.code
      identifier.make_global!
      return v, defs[:__MODULE__]
    end
    defs[:__IMPORTS__].each do |imported_defs|
      if v = imported_defs[var]
        identifier.code = v.code
        identifier.make_global!
        return v, imported_defs[:__MODULE__]
      end
    end
    nil
  end

  def self.def_def(current, identifier)
    identifier.make_global!
    existing, _ = get_module_def?(current, identifier)
    if existing
      identifier.code = existing.code
      return
    end
    if identifier.no_module?
      identifier.add_module! current.whole
    end
    identifier.code = @@next_code
    @@next_code += 1
    defs = @@defs[identifier.module_part] ||= blank_def(identifier)
    defs[identifier.var_part] = identifier
  end

  def self.blank_def(mod_iden)
    unless @@defs['Kernel']
      @@kernel = @@defs['Kernel'] = {
        __MODULE__: Identifier.new('Kernel', nil),
        __ALIAS__: Hash.new,
        __IMPORTS__: []
      }
      if mod_iden.whole == 'Kernel'
        return @@defs['Kernel']
      end
    end
    {
      __MODULE__: mod_iden,
      __ALIAS__: Hash.new,
      __IMPORTS__: [@@kernel]
    }
  end

  def self.define_module(identifier)
    identifier.make_global!
    mod = identifier.whole
    @@defs[mod] ||= blank_def(identifier)
    if identifier.code.nil?
      identifier.code = @@next_code
      @@next_code += 1
    end
    identifier
  end

  def self.get_next_code
    id = @@next_code
    @@next_code += 1
    id
  end

  # Get the code for a certain statically-known method
  def self.[](mod, method)
    mod = mod.to_s
    method = method.to_s
    methods = @@defs[mod]
    raise "Can't find method #{mod}.#{method}" unless methods
    iden = methods[method]
    raise "Can't find method #{mod}.#{method}" unless iden
    iden.code
  end
end

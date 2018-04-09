require 'set'

class DAG
  def initialize
    @edges = Hash.new
  end

  def add(from, to)
    @edges[from] ||= Set.new
    raise "Cyclic dependency! #{from} #{to}" if route_to(to, from)
    @edges[from] << to
  end

  def route_to(src, dest)
    visited = Set.new
    to_visit = [src]
    while n = to_visit.shift
      return true if dest == n
      next if visited.include? n
      visited << n
      deps = @edges[n]
      if deps
        to_visit += deps.to_a
      end
    end
    false
  end

  def dependent_order(start)
    res = []
    done = Set.new
    to_visit = [start]
    while n = to_visit.shift
      next if done.include? n
      done << n
      res << n
      deps = @edges[n]
      if deps
        to_visit += deps.to_a
      end
    end
    res
  end
end

class DAGOrder
  def self.order(module_contents, dependencies)
    dag = DAG.new
    dependencies.each do |mod, uses|
      uses.each do |other|
        dag.add(mod, other)
      end
    end
    result = []
    dag.dependent_order("Main").reverse.each do |mod|
      if contents = module_contents[mod]
        result += contents
      end
    end
    result
  end
end

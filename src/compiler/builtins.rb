require_relative './builtins_dsl'

defmodule :IO do
  defn :puts
  defn :gets
end

defmodule :Kernel do
  defn :type
  defn :<
  defn :-
  defn :*
  defn :conj
end

deftype :Channel do
  defn :new
  defn :send
  defn :receive
  # defimpl :conj
end

defmodule :Sequence do
  defprotocol :cons
  defprotocol :conj
  defprotocol :head
  defprotocol :tail
end

defmodule :Enumerable do
  defprotocol :reduce
end

deftype :List do
  defn :new
  defimpl :conj, of: [:Sequence, :conj]
  defimpl :head, of: [:Sequence, :head]
  defimpl :tail, of: [:Sequence, :tail]
end

deftype :Vector do
  defimpl :conj, of: [:Sequence, :conj]
  defimpl :head, of: [:Sequence, :head]
  defimpl :tail, of: [:Sequence, :tail]
end

Def.assign_codes

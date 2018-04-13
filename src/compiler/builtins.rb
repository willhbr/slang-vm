require_relative './builtins_dsl'

defmodule :Printable do
  defprotocol :"->string"
end

defmodule :IO do
  defn :puts
  defn :gets
end

defmodule :Kernel do
  defn :type
  defn :<
  defn :-
  defn :*
end

deftype :Int do
end

deftype :String do
end

deftype :Channel do
  defn :new
  defn :send
  defn :receive
  # defimpl :conj
end

defmodule :Sequence do
  defprotocol :conj
  defprotocol :head
  defprotocol :tail
end

defmodule :Enumerable do
  defprotocol :reduce
end

deftype :List do
  defimpls :Sequence
  defn :new
end

deftype :Vector do
  defimpls :Sequence
end

Def.assign_codes

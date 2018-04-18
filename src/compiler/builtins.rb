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
  defn :+
  defn :/
  defn :*
end

defmodule :Access do
  defprotocol :get
  # defprotocol :set
end

defmodule :Equatable do
  defprotocol :"="
end

deftype :Int, '*big.Int' do
  defimpls :Printable
end

deftype :String, 'string' do
  defimpls :Printable
  defimpls :Access
  defimpls :Equatable
end

deftype :Atom do
  defimpls :Printable
  defn :value
end

defmodule :File do
  defn :read
end

deftype :Channel, 'chan Value' do
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

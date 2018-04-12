require_relative './objects'

class Builtins
  MODULES = {
    IO: [
      :puts,
      :gets
    ],
    Kernel: [
      :type,
      :<,
      :-,
      :*,
      :conj
    ],
    Channel: [
      :new,
      :send,
      :receive
    ],
    Enumerable: [
      :reduce
    ]
  }

  PROTOCOL_METHODS = {
    Enumerable: [
      :reduce
    ]
  }
end


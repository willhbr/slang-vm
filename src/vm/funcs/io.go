package funcs

import (
	"../ds"
	"../vm"
	"bufio"
	"fmt"
	"os"
)

func IO__puts(co *vm.Coroutine, arguments ...ds.Value) ds.Value {
	for i := range arguments {
		arg := arguments[i]
		// TODO Call a real method to turn things into a string
		switch arg.(type) {
		case ds.Atom:
			fmt.Print(co.Program.Strings[int(arg.(ds.Atom))])
		default:
			fmt.Print(arguments[i])
		}
	}
	fmt.Println()
	return ds.Nil
}

func IO__gets(co *vm.Coroutine, arguments ...ds.Value) ds.Value {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text
}

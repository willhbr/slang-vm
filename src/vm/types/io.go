package types

import (
	"bufio"
	"fmt"
	"os"
)

func IO__puts(_ *Program, arguments ...Value) (Value, error) {
	for i := range arguments {
		fmt.Print(arguments[i])
	}
	fmt.Println()
	return Nil, nil
}

func IO__gets(_ *Program, arguments ...Value) (Value, error) {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text, nil
}

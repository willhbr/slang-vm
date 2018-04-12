package types

import (
	"bufio"
	"fmt"
	"os"
)

func IO__puts(arguments ...Value) Value {
	for i := range arguments {
		fmt.Print(arguments[i])
	}
	fmt.Println()
	return Nil
}

func IO__gets(arguments ...Value) Value {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text
}

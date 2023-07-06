package utils

import (
	"fmt"
	"os"
)

type CommandUtils struct {
}

func (CommandUtils) PressAnyKeyToContinue() {
	fmt.Println("\nPress enter key to exit.")
	var input string
	fmt.Scanln(&input)
	os.Exit(0)
}

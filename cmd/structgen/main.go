package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	var err error

	if len(os.Args) <= 1 {
		fmt.Println("structgen: No argument specified")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "theme":
		err = generateTheme()

	case "keybindings":
		err = generateKeybindings()

	default:
		err = errors.New("invalid argument")
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

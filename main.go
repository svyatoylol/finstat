package main

import (
	"fin-stat/commands"
	"fmt"
)

func main() {
	var command string
	fmt.Println("Enter a command:")
	fmt.Scanln(&command)

	commands.ExecuteCommand(command)
}

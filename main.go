package main

import (
	"fin-stat/commands"
	"fin-stat/transaction"
	"fmt"
)

func main() {
	var command string
	fmt.Println("Enter a command:")
	fmt.Scanln(&command)

	repo := transaction.NewMemoryRepository()
	svc := transaction.NewService(repo)

	commands.ExecuteCommand(command, svc)
}

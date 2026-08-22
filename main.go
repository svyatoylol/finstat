package main

import (
	"bufio"
	"fin-stat/commands"
	"fin-stat/transaction"
	"fmt"
	"os"
)

func main() {
	repo := transaction.NewMemoryRepository()
	svc := transaction.NewService(repo)
	sess := commands.NewSession(0)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter commands (type 'exit' to quit):")

	for scanner.Scan() {
		cmd := commands.ParseCommand(scanner.Text())
		if cmd.Name == "exit" {
			break
		}
		commands.ExecuteCommand(cmd, svc, sess)
		fmt.Println("Enter a command:")
	}
}

package commands

import (
	"context"
	"errors"
	"fin-stat/transaction"
	"fmt"
	"log"
)

func ExecuteCommand(ctx context.Context, cmd ParsedCommand, svc *transaction.Service, sess *Session) {
	if cmd.Name == "Login" {
		handleLogin(cmd.Args, sess)
		return
	}

	if !sess.loggedIn {
		fmt.Println("You must be logged in to execute this command.")
		return
	}

	switch cmd.Name {
	case "Add_Incoming":
		handleAddTransaction(ctx, cmd.Args, transaction.Incoming, svc, sess)
	case "Add_Outgoing":
		handleAddTransaction(ctx, cmd.Args, transaction.Outgoing, svc, sess)
	case "Check_Balance":
		handleCheckBalance(ctx, svc, sess)
	case "All_By_Type":
		handleListByType(ctx, cmd.Args, svc, sess)
	case "All_By_Category":
		handleListByCategory(ctx, cmd.Args, svc, sess)
	case "Delete_Transaction":
		handleDeleteTransaction(ctx, cmd.Args, svc, sess)
	case "Edit_Transaction":
		handleEditTransaction(ctx, cmd.Args, svc, sess)
	default:
		fmt.Println("Unknown command:", cmd.Name)
	}
}

// printError показывает пользователю понятный текст для ожидаемых ошибок,
// а внутренние (БД и т.п.) логирует и скрывает за общим сообщением.
func printError(err error) {
	switch {
	case errors.Is(err, transaction.ErrNotFound):
		fmt.Println("transaction not found")
	case errors.Is(err, transaction.ErrInvalidAmount):
		fmt.Println("amount must be positive")
	case errors.Is(err, transaction.ErrInvalidType):
		fmt.Println("invalid transaction type")
	case errors.Is(err, context.DeadlineExceeded):
		fmt.Println("request timed out, try again")
	default:
		log.Printf("internal error: %v", err)
		fmt.Println("something went wrong, try again later")
	}
}

package commands

import (
	"fin-stat/transaction"
	"fmt"
	"strconv"
)

func HandeleAddTransaction(cmd ParsedCommand, txType transaction.TransactionType, svc *transaction.Service) {
	if len(cmd.Args) < 2 {
		fmt.Println("Invalid number of arguments for Add_Transaction command")
		fmt.Println("usage: <amount> <category>")
		return
	}

	amount, err := strconv.ParseFloat(cmd.Args[0], 64)
	if err != nil {
		fmt.Println("Invalid amount:", cmd.Args[0])
		return
	}

	t := transaction.Transaction{
		Type:     txType,
		Amount:   amount,
		Category: cmd.Args[1],
	}

	saved, err := svc.AddTransaction(t)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("Added transaction #", saved.ID)
}

func handleDeleteTransaction(args []string, svc *transaction.Service) {
	if len(args) < 1 {
		fmt.Println("usage: Delete_Transaction <id>")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Invalid id:", args[0])
		return
	}

	if err := svc.DeleteTransaction(id); err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("Deleted transaction #", id)
}

func handleEditTransaction(args []string, svc *transaction.Service) {
	if len(args) < 3 {
		fmt.Println("usage: Edit_Transaction <id> <amount> <category>")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("invalid id:", args[0])
		return
	}

	amount, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		fmt.Println("invalid amount:", args[1])
		return
	}

	existing, ok := svc.GetByID(id)
	if !ok {
		fmt.Println("transaction not found:", id)
		return
	}

	existing.Amount = amount
	existing.Category = args[2]

	updated, err := svc.UpdateTransaction(existing)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("Updated transaction #", updated.ID)
}

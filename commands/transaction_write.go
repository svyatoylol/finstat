package commands

import (
	"context"
	"fin-stat/transaction"
	"fmt"
	"strconv"
)

func handleAddTransaction(ctx context.Context, args []string, txType transaction.TransactionType, svc *transaction.Service, sess *Session) {
	if len(args) < 2 {
		fmt.Println("usage: Add_Incoming|Add_Outgoing <amount> <category>")
		return
	}

	amount, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		fmt.Println("invalid amount:", args[0])
		return
	}

	saved, err := svc.AddTransaction(ctx, transaction.Transaction{
		UserID:   sess.UserID,
		Type:     txType,
		Amount:   amount,
		Category: args[1],
	})
	if err != nil {
		printError(err)
		return
	}

	fmt.Printf("Added transaction #%d\n", saved.ID)
}

func handleDeleteTransaction(ctx context.Context, args []string, svc *transaction.Service, sess *Session) {
	if len(args) < 1 {
		fmt.Println("usage: Delete_Transaction <id>")
		return
	}

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		fmt.Println("invalid id:", args[0])
		return
	}

	if err := svc.DeleteTransaction(ctx, sess.UserID, id); err != nil {
		printError(err)
		return
	}

	fmt.Printf("Deleted transaction #%d\n", id)
}

func handleEditTransaction(ctx context.Context, args []string, svc *transaction.Service, sess *Session) {
	if len(args) < 3 {
		fmt.Println("usage: Edit_Transaction <id> <amount> <category>")
		return
	}

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		fmt.Println("invalid id:", args[0])
		return
	}

	amount, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		fmt.Println("invalid amount:", args[1])
		return
	}

	existing, err := svc.GetByID(ctx, sess.UserID, id)
	if err != nil {
		printError(err)
		return
	}

	existing.Amount = amount
	existing.Category = args[2]

	updated, err := svc.UpdateTransaction(ctx, sess.UserID, existing)
	if err != nil {
		printError(err)
		return
	}

	fmt.Printf("Updated transaction #%d\n", updated.ID)
}

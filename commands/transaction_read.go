package commands

import (
	"context"
	"fin-stat/transaction"
	"fmt"
	"strconv"
)

func handleLogin(args []string, sess *Session) {
	if len(args) < 1 {
		fmt.Println("usage: Login <user_id>")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("invalid user id:", args[0])
		return
	}

	sess.UserID = id
	sess.loggedIn = true
	fmt.Printf("Logged in as user %d\n", id)
}

func handleCheckBalance(ctx context.Context, svc *transaction.Service, sess *Session) {
	balance, err := svc.GetBalance(ctx, sess.UserID)
	if err != nil {
		printError(err)
		return
	}
	fmt.Printf("Balance: %.2f\n", balance)
}

func handleListByType(ctx context.Context, args []string, svc *transaction.Service, sess *Session) {
	if len(args) < 1 {
		fmt.Println("usage: All_By_Type <incoming|outgoing>")
		return
	}

	txType := transaction.TransactionType(args[0])
	if txType != transaction.Incoming && txType != transaction.Outgoing {
		fmt.Println("invalid type:", args[0])
		return
	}

	txs, err := svc.GetAllByType(ctx, sess.UserID, txType)
	if err != nil {
		printError(err)
		return
	}
	printTransactions(txs)
}

func handleListByCategory(ctx context.Context, args []string, svc *transaction.Service, sess *Session) {
	if len(args) < 1 {
		fmt.Println("usage: All_By_Category <category>")
		return
	}

	txs, err := svc.GetAllByCategory(ctx, sess.UserID, args[0])
	if err != nil {
		printError(err)
		return
	}
	printTransactions(txs)
}

func printTransactions(txs []transaction.Transaction) {
	if len(txs) == 0 {
		fmt.Println("no transactions found")
		return
	}
	for _, t := range txs {
		fmt.Printf("#%d  %-8s  %.2f  %s\n", t.ID, t.Type, t.Amount, t.Category)
	}
}

package commands

import (
	"fin-stat/transaction"
	"fmt"
)

func handleListByType(args []string, svc *transaction.Service) {
	if len(args) < 1 {
		fmt.Println("usage: List_By_Type <incoming|outgoing>")
		return
	}

	txType := transaction.TransactionType(args[0])
	if txType != transaction.Incoming && txType != transaction.Outgoing {
		fmt.Println("invalid type:", args[0])
		return
	}

	printTransactions(svc.GetAllByType(txType))
}

func handleListByCategory(args []string, svc *transaction.Service) {
	if len(args) < 1 {
		fmt.Println("usage: List_By_Category <category>")
		return
	}

	txCategory := args[0]
	printTransactions(svc.GetAllByCategory(txCategory))
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

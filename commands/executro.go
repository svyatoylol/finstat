package commands

import (
	"fin-stat/transaction"
	"fmt"
	"strconv"
)

func ExecuteCommand(cmd ParsedCommand, svc *transaction.Service) {
	switch cmd.Name {
	case "Add_Incoming":
		fmt.Println("Executing Add_Incoming command...")
		HandeleAddTransaction(cmd, transaction.Incoming, svc)
	case "Add_Outgoing":
		fmt.Println("Executing Add_Outgoing command...")
		HandeleAddTransaction(cmd, transaction.Outgoing, svc)
	case "Check_Balance":
		fmt.Println(svc.GetBalance())
	case "All_By_Type":
		handleListByType(cmd.Args, svc)
	case "All_By_Category":
		handleListByCategory(cmd.Args, svc)
	default:
		fmt.Println("Unknown command:", cmd.Name)
	}
}

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

	t := transaction.Transaction{
		ID:       id,
		Amount:   amount,
		Category: args[2],
	}

	updated, err := svc.UpdateTransaction(t)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("Updated transaction #", updated.ID)
}

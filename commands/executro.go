package commands

import (
	"fin-stat/transaction"
	"fmt"
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
	case "Delete_Transaction":
		handleDeleteTransaction(cmd.Args, svc)
	case "Edit_Transaction":
		handleEditTransaction(cmd.Args, svc)
	default:
		fmt.Println("Unknown command:", cmd.Name)
	}
}

package commands

import (
	"fin-stat/transaction"
	"fmt"
)

func ExecuteCommand(cmd ParsedCommand, svc *transaction.Service, sess *Session) {
	if cmd.Name == "Login" {
		handleLoggin(cmd.Args, sess)
		return
	}

	if !sess.loggedIn {
		fmt.Println("You must be logged in to execute this command.")
		return
	}

	switch cmd.Name {
	case "Add_Incoming":
		fmt.Println("Executing Add_Incoming command...")
		HandeleAddTransaction(cmd, transaction.Incoming, svc)
	case "Add_Outgoing":
		fmt.Println("Executing Add_Outgoing command...")
		HandeleAddTransaction(cmd, transaction.Outgoing, svc)
	case "Check_Balance":
		fmt.Println(svc.GetBalance(sess.UserID))
	case "All_By_Type":
		handleListByType(cmd.Args, svc, sess)
	case "All_By_Category":
		handleListByCategory(cmd.Args, svc, sess)
	case "Delete_Transaction":
		handleDeleteTransaction(cmd.Args, svc)
	case "Edit_Transaction":
		handleEditTransaction(cmd.Args, svc)
	default:
		fmt.Println("Unknown command:", cmd.Name)
	}
}

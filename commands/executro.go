package commands

import (
	"fin-stat/transaction"
	"fmt"
)

func ExecuteCommand(command string, svc *transaction.Service) {
	switch command {
	case "Add_Incoming":
		fmt.Println("Executing Add_Incoming command...")
		// Add your logic for Add_Incoming command here
	case "Add_Transaction":
		fmt.Println("Executing Add_Transaction command...")
		// Add your logic for Add_Transaction command here
	case "Check_Balance":
		fmt.Println("Executing Check_Balance command...")
		// Add your logic for Check_Balance command here
	default:
		fmt.Println("Unknown command:", command)
	}
}

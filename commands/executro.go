package commands

import "fmt"

func ExecuteCommand(command string) {
	switch command {
	case "Add_Incoming":
		fmt.Println("Executing Add_Incoming command...")
		// Add your logic for Add_Incoming command here
	case "Add_Outgoing":
		fmt.Println("Executing Add_Outgoing command...")
		// Add your logic for Add_Outgoing command here
	case "Add_Transfer":
		fmt.Println("Executing Add_Transfer command...")
		// Add your logic for Add_Transfer command here
	default:
		fmt.Println("Unknown command:", command)
	}
}

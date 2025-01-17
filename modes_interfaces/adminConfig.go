package modes_interfaces

import (
	"Nospresso/coffee_machines"
	"fmt"
)

const (
	AdminRestock    = "1"
	AdminUpdatePIN  = "2"
	AdminReturnMenu = "3"
	attempts        = 3
)

func ConfigAdmin(machines []coffee_machines.Machine) {
	if len(machines) == 0 {
		fmt.Println("No machines available.")
		return
	}
	fmt.Println("Select a machine by ID (1 to", len(machines), ")")
	machineID := coffee_machines.GetValidatedNumber("> ", 1, len(machines))
	currentMachine := &machines[machineID-1]

	if !coffee_machines.ValidatePin(currentMachine, attempts) {
		fmt.Println("Too many failed attempts. Exiting admin mode.")
		return
	}
	adminChoice := coffee_machines.GetValidatedInput("Access granted.\n1) Restock\n2) Update PIN\n3) Return", []string{"1", "2", "3"})
	switch adminChoice {
	case AdminRestock:
		coffee_machines.HandleRestock(currentMachine)
	case AdminUpdatePIN:
		coffee_machines.UpdatePin(currentMachine)
	case AdminReturnMenu:
		fmt.Println("Returning to main menu.")
	}
}

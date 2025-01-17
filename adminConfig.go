package main

import "fmt"

func configAdmin(machines []Machine) {
	if len(machines) == 0 {
		fmt.Println("No machines available.")
		return
	}
	fmt.Println("Select a machine by ID (1 to", len(machines), ")")
	machineID := GetValidatedNumber("> ", 1, len(machines))
	currentMachine := &machines[machineID-1]
	const attempts = 3
	if !ValidatePin(currentMachine, attempts) {
		fmt.Println("Too many failed attempts. Exiting admin mode.")
		return
	}
	adminChoice := GetValidatedInput("Access granted.\n1) Restock\n2) Update PIN\n3) Return", []string{"1", "2", "3"})
	switch adminChoice {
	case AdminRestock:
		handleRestock(currentMachine)
	case AdminUpdatePIN:
		UpdatePin(currentMachine)
	case AdminReturnMenu:
		fmt.Println("Returning to main menu.")
	}
}

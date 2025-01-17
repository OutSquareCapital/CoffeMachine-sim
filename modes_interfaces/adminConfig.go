package modes_interfaces

import (
	"Nospresso/coffee_machines"
	"Nospresso/validations"
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	AdminRestock    = "1"
	AdminUpdatePIN  = "2"
	AdminReturnMenu = "3"
	attempts        = 3
)

func ValidatePin(machine *coffee_machines.Machine, attempts int) bool {
	reader := bufio.NewReader(os.Stdin)
	for attempts > 0 {
		fmt.Println("Enter PIN:")
		pin, _ := reader.ReadString('\n')
		pin = strings.TrimSpace(pin)
		if pin == machine.Pincode {
			return true
		}
		attempts--
		if attempts > 0 {
			fmt.Printf("Incorrect PIN. %d attempts remaining.\n", attempts)
		}
	}
	return false
}

func UpdatePin(machine *coffee_machines.Machine) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter new 6-digit PIN:")
	for {
		pin, _ := reader.ReadString('\n')
		pin = strings.TrimSpace(pin)
		if len(pin) == 6 && validations.IsNumeric(pin) {
			machine.Pincode = pin
			fmt.Println("PIN updated successfully.")
			break
		}
		fmt.Println("Invalid PIN. Enter a new 6-digit PIN:")
	}
}

func ConfigAdmin(machines []coffee_machines.Machine) {
	if len(machines) == 0 {
		fmt.Println("No machines available.")
		return
	}
	fmt.Println("Select a machine by ID (1 to", len(machines), ")")
	machineID := validations.GetValidatedNumber("> ", 1, len(machines))
	currentMachine := &machines[machineID-1]

	if !ValidatePin(currentMachine, attempts) {
		fmt.Println("Too many failed attempts. Exiting admin mode.")
		return
	}
	adminChoice := validations.GetValidatedInput("Access granted.\n1) Restock\n2) Update PIN\n3) Return", []string{"1", "2", "3"})
	switch adminChoice {
	case AdminRestock:
		coffee_machines.HandleRestock(currentMachine)
	case AdminUpdatePIN:
		UpdatePin(currentMachine)
	case AdminReturnMenu:
		fmt.Println("Returning to main menu.")
	}
}

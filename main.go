package main

import (
	"Nospresso/coffee_machines"
	"Nospresso/modes_interfaces"
	"Nospresso/validations"
	"fmt"
)

const (
	filename   = "machines.csv"
	ModeClient = "1"
	ModeAdmin  = "2"
	ModeExit   = "3"
)

func main() {
	machines, err := coffee_machines.InitializeMachines(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	coffee_machines.DisplayMachines(machines)
	beverages := coffee_machines.InitializeBeverages()

	for {
		fmt.Println("Nospresso Cafe")
		fmt.Printf("%s) Client\n%s) Admin\n%s) Exit\n", ModeClient, ModeAdmin, ModeExit)
		mode := validations.GetValidatedInput("> ", []string{
			ModeClient,
			ModeAdmin,
			ModeExit})
		switch mode {
		case ModeClient:
			modes_interfaces.ServeClient(machines, beverages)
		case ModeAdmin:
			modes_interfaces.ConfigAdmin(machines)
		case ModeExit:
			fmt.Println("Saving machines to file and exiting...")
			if err := coffee_machines.SaveCSV(filename, machines); err != nil {
				fmt.Println("Error saving machines:", err)
			}
			return
		}
	}
}

package main

import (
	"fmt"
)

const (
	filename        = "machines.csv"
	ModeClient      = "1"
	ModeAdmin       = "2"
	ModeExit        = "3"
	AdminRestock    = "1"
	AdminUpdatePIN  = "2"
	AdminReturnMenu = "3"
)

func main() {
	machines, err := LoadCSV(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	DisplayMachines(machines)
	for {
		fmt.Println("Nospresso Cafe")
		fmt.Println("Select mode:")
		fmt.Printf("%s) Client\n%s) Admin\n%s) Exit\n", ModeClient, ModeAdmin, ModeExit)
		mode := GetValidatedInput("> ", []string{ModeClient, ModeAdmin, ModeExit})
		switch mode {
		case ModeClient:
			serveClient(machines)
		case ModeAdmin:
			configAdmin(machines)
		case ModeExit:
			fmt.Println("Saving machines to file and exiting...")
			if err := SaveCSV(filename, machines); err != nil {
				fmt.Println("Error saving machines:", err)
			}
			return
		}
	}
}

package coffee_machines

import (
	"Nospresso/validations"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Machine struct {
	ID        int
	Pincode   string
	Inventory Inventory
	Beverages []Beverage
}

func InitializeMachines(filename string) ([]Machine, error) {
	machines, err := LoadCSV(filename)
	if err != nil {
		return nil, err
	}

	beverages := InitializeBeverages()
	for i := range machines {
		machines[i].Beverages = beverages
	}
	return machines, nil
}

func SelectMachine(machines []Machine) *Machine {
	fmt.Println("Select a machine by ID (1 to", len(machines), ")")
	machineID := validations.GetValidatedNumber("> ", 1, len(machines))
	return &machines[machineID-1]
}

func LoadCSV(filename string) ([]Machine, error) {
	var machines []Machine
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("error: CSV file not found, exiting program")
	}
	defer file.Close()
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New("invalid CSV format")
	}
	for id, line := range lines[1:] {
		if len(line) != 4 {
			return nil, errors.New("invalid data in CSV")
		}
		milk, _ := strconv.Atoi(line[1])
		sugar, _ := strconv.Atoi(line[2])
		coffee, _ := strconv.Atoi(line[3])
		machines = append(machines, Machine{
			ID:      id + 1,
			Pincode: line[0],
			Inventory: Inventory{
				Coffee: Ingredient{Name: "Coffee", Quantity: coffee},
				Milk:   Ingredient{Name: "Milk", Quantity: milk},
				Sugar:  Ingredient{Name: "Sugar", Quantity: sugar},
			},
		})
	}
	return machines, nil
}

func SaveCSV(filename string, machines []Machine) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"PINCODE", "MILK", "SUGAR", "COFFEE"})
	for _, m := range machines {
		writer.Write([]string{
			m.Pincode,
			strconv.Itoa(m.Inventory.Milk.Quantity),
			strconv.Itoa(m.Inventory.Sugar.Quantity),
			strconv.Itoa(m.Inventory.Coffee.Quantity),
		})
	}
	return nil
}

func DisplayMachines(machines []Machine) {
	if len(machines) == 0 {
		fmt.Println("No machines loaded.")
		return
	}
	for _, machine := range machines {
		fmt.Printf("Machine %d loaded:\n", machine.ID)
		fmt.Printf("ID: %d\nPincode: %s\n", machine.ID, machine.Pincode)
		machine.Inventory.DisplayStock()
		fmt.Println()
	}
}

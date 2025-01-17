package coffee_machines

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Machine struct {
	ID        int
	Pincode   string
	Inventory Inventory
}

func SelectMachine(machines []Machine) *Machine {
	fmt.Println("Select a machine by ID (1 to", len(machines), ")")
	machineID := GetValidatedNumber("> ", 1, len(machines))
	return &machines[machineID-1]
}

func (m *Machine) RemoveIngredient(ingredient string, amount int) bool {
	switch ingredient {
	case "milk":
		if m.Inventory.Milk.Quantity >= amount {
			m.Inventory.Milk.Quantity -= amount
			return true
		}
	case "sugar":
		if m.Inventory.Sugar.Quantity >= amount {
			m.Inventory.Sugar.Quantity -= amount
			return true
		}
	case "coffee":
		if m.Inventory.Coffee.Quantity >= amount {
			m.Inventory.Coffee.Quantity -= amount
			return true
		}
	}
	return false
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

func ValidatePin(machine *Machine, attempts int) bool {
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

func UpdatePin(machine *Machine) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter new 6-digit PIN:")
	for {
		pin, _ := reader.ReadString('\n')
		pin = strings.TrimSpace(pin)
		if len(pin) == 6 && isNumeric(pin) {
			machine.Pincode = pin
			fmt.Println("PIN updated successfully.")
			break
		}
		fmt.Println("Invalid PIN. Enter a new 6-digit PIN:")
	}
}

package main

import (
	"encoding/csv"
	"errors"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Ingredient struct {
	Name         string
	Quantity     int
	PricePerUnit float64
}

type Inventory struct {
	Coffee Ingredient
	Milk   Ingredient
	Sugar  Ingredient
}

type Beverage struct {
	Name      string
	Coffee    int
	Milk      int
	Sugar     int
	BasePrice float64
}

var Beverages = map[string]Beverage{
	"Espresso":       {Name: "Espresso", Coffee: 8, Milk: 0, Sugar: 0, BasePrice: 2.00},
	"Cappuccino":     {Name: "Cappuccino", Coffee: 6, Milk: 100, Sugar: 0, BasePrice: 2.50},
	"Latte (Small)":  {Name: "Latte (Small)", Coffee: 6, Milk: 120, Sugar: 0, BasePrice: 2.70},
	"Latte (Medium)": {Name: "Latte (Medium)", Coffee: 8, Milk: 150, Sugar: 0, BasePrice: 3.20},
	"Latte (Large)":  {Name: "Latte (Large)", Coffee: 12, Milk: 200, Sugar: 0, BasePrice: 3.70},
}

func (inv *Inventory) VerifyStock(coffee, sugar, milk int) error {
	if inv.Coffee.Quantity < coffee {
		return errors.New("insufficient coffee stock")
	}
	if inv.Sugar.Quantity < sugar {
		return errors.New("insufficient sugar stock")
	}
	if inv.Milk.Quantity < milk {
		return errors.New("insufficient milk stock")
	}
	return nil
}

func (inv *Inventory) UpdateStock(coffee, sugar, milk int) {
	inv.Coffee.Quantity -= coffee
	inv.Sugar.Quantity -= sugar
	inv.Milk.Quantity -= milk
}

type Machine struct {
	ID        int
	Pincode   string
	Inventory Inventory
}

func LoadCSV(filename string) ([]Machine, error) {
	var machines []Machine
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("CSV file not found")
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

func GenerateTwintCode() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var code strings.Builder
	for i := 0; i < 5; i++ {
		code.WriteByte(chars[rand.Intn(len(chars))])
	}
	return code.String()
}

func ValidatePin(machine *Machine, pin string) bool {
	return pin == machine.Pincode
}

func UpdatePin(machine *Machine, newPin string) error {
	if len(newPin) != 6 || !isNumeric(newPin) {
		return errors.New("invalid PIN format")
	}
	machine.Pincode = newPin
	return nil
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func getMachineIDs(machines []Machine) []string {
	var ids []string
	for _, m := range machines {
		ids = append(ids, strconv.Itoa(m.ID))
	}
	return ids
}

func calculatePrice(beverage Beverage, sugarOption, milkOption string) (float64, error) {
	var sugarPrice, milkPrice float64

	// Ajout du prix du sucre
	switch sugarOption {
	case "Light (5g)":
		sugarPrice = 0.10
	case "Medium (10g)":
		sugarPrice = 0.20
	case "Heavy (15g)":
		sugarPrice = 0.30
	}

	// Ajout du prix du lait supplémentaire
	switch milkOption {
	case "1 dose":
		milkPrice = 0.05
	case "2 doses":
		milkPrice = 0.10
	case "3 doses":
		milkPrice = 0.15
	}

	return beverage.BasePrice + sugarPrice + milkPrice, nil
}

func getBeverageRequirements(beverageName, sugarOption, milkOption string) (Beverage, int, error) {
	beverage, exists := Beverages[beverageName]
	if !exists {
		return Beverage{}, 0, errors.New("boisson inconnue")
	}

	// Ajout du sucre en fonction de l'option choisie
	switch sugarOption {
	case "Light (5g)":
		beverage.Sugar = 5
	case "Medium (10g)":
		beverage.Sugar = 10
	case "Heavy (15g)":
		beverage.Sugar = 15
	}

	// Ajout du lait supplémentaire
	switch milkOption {
	case "1 dose":
		beverage.Milk += 50
	case "2 doses":
		beverage.Milk += 100
	case "3 doses":
		beverage.Milk += 150
	}

	return beverage, beverage.Milk, nil
}

package coffee_machines

import "fmt"

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

func (inv *Inventory) DisplayStock() {
	fmt.Printf("Current stock levels:\n")
	fmt.Printf("- %s: %dg\n", inv.Coffee.Name, inv.Coffee.Quantity)
	fmt.Printf("- %s: %.3fL\n", inv.Milk.Name, float64(inv.Milk.Quantity)/1000.0)
	fmt.Printf("- %s: %dg\n", inv.Sugar.Name, inv.Sugar.Quantity)
}

func (inv *Inventory) UpdateStock(coffee, sugar, milk int) {
	inv.Coffee.Quantity -= coffee
	inv.Sugar.Quantity -= sugar
	inv.Milk.Quantity -= milk
}

func HandleRestock(machine *Machine) {
	machine.Inventory.DisplayStock()
	coffee := GetValidatedQuantity("Enter coffee powder (g) to add > ")
	sugar := GetValidatedQuantity("Enter sugar (g) to add > ")
	milk := GetValidatedQuantity("Enter milk (L) to add > ")
	machine.Inventory.Coffee.Quantity += coffee
	machine.Inventory.Sugar.Quantity += sugar
	machine.Inventory.Milk.Quantity += milk * 1000
	fmt.Println("Stocks updated successfully.")
}

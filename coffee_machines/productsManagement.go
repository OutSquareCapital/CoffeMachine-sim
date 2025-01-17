package coffee_machines

import (
	"Nospresso/validations"
	"fmt"
)

const (
	espresso   = "Espresso"
	cappuccino = "Cappuccino"
	latte      = "Latte"
)

type Ingredient struct {
	Name         string
	Quantity     int
	PricePerUnit float64
}

type Beverage struct {
	Name        string
	BasePrice   float64
	Ingredients map[string]int
}

type Inventory struct {
	Coffee Ingredient
	Milk   Ingredient
	Sugar  Ingredient
}

func InitializeBeverages() []Beverage {
	return []Beverage{
		{
			Name:      espresso,
			BasePrice: 2.00,
			Ingredients: map[string]int{
				"coffee": 8,
			},
		},
		{
			Name:      cappuccino,
			BasePrice: 2.50,
			Ingredients: map[string]int{
				"coffee": 6,
				"milk":   100,
			},
		},
		{
			Name:      latte,
			BasePrice: 2.70,
			Ingredients: map[string]int{
				"coffee": 6,
				"milk":   120,
			},
		},
	}
}

func (inv *Inventory) CalculatePrice(beverage Beverage) float64 {
	totalPrice := beverage.BasePrice
	for ingredientName, quantity := range beverage.Ingredients {
		var pricePerUnit float64
		switch ingredientName {
		case "coffee":
			pricePerUnit = inv.Coffee.PricePerUnit
		case "milk":
			pricePerUnit = inv.Milk.PricePerUnit
		case "sugar":
			pricePerUnit = inv.Sugar.PricePerUnit
		}
		totalPrice += float64(quantity) * pricePerUnit
	}
	return totalPrice
}

func (inv *Inventory) DisplayStock() {
	fmt.Printf("Current stock levels:\n")
	fmt.Printf("- %s: %dg\n", inv.Coffee.Name, inv.Coffee.Quantity)
	fmt.Printf("- %s: %.3fL\n", inv.Milk.Name, float64(inv.Milk.Quantity)/1000.0)
	fmt.Printf("- %s: %dg\n", inv.Sugar.Name, inv.Sugar.Quantity)
}

func (inv *Inventory) VerifyStock(beverage Beverage) bool {
	for ingredientName, quantity := range beverage.Ingredients {
		var stock int
		switch ingredientName {
		case "coffee":
			stock = inv.Coffee.Quantity
		case "milk":
			stock = inv.Milk.Quantity
		case "sugar":
			stock = inv.Sugar.Quantity
		}
		if stock < quantity {
			fmt.Printf("Error: Insufficient %s to prepare %s.\n", ingredientName, beverage.Name)
			return false
		}
	}
	return true
}

func (inv *Inventory) UpdateStock(beverage Beverage) {
	for ingredientName, quantity := range beverage.Ingredients {
		switch ingredientName {
		case "coffee":
			inv.Coffee.Quantity -= quantity
		case "milk":
			inv.Milk.Quantity -= quantity
		case "sugar":
			inv.Sugar.Quantity -= quantity
		}
	}
}

func HandleRestock(machine *Machine) {
	machine.Inventory.DisplayStock()
	coffee := validations.GetValidatedQuantity("Enter coffee powder (g) to add > ")
	sugar := validations.GetValidatedQuantity("Enter sugar (g) to add > ")
	milk := validations.GetValidatedQuantity("Enter milk (L) to add > ")
	machine.Inventory.Coffee.Quantity += coffee
	machine.Inventory.Sugar.Quantity += sugar
	machine.Inventory.Milk.Quantity += milk * 1000
	fmt.Println("Stocks updated successfully.")
}

package modes_interfaces

import (
	"Nospresso/coffee_machines"
	"Nospresso/validations"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func chooseBeverage(beverages []coffee_machines.Beverage) coffee_machines.Beverage {
	fmt.Println("Please select your beverage:")
	for i, bev := range beverages {
		fmt.Printf("%d) %s - CHF %.2f\n", i+1, bev.Name, bev.BasePrice)
	}
	selected := validations.GetValidatedNumber("> ", 1, len(beverages))
	return beverages[selected-1]
}

func ServeClient(machines []coffee_machines.Machine, beverages []coffee_machines.Beverage) {
	if len(machines) == 0 {
		fmt.Println("No machines available.")
		return
	}
	machine := coffee_machines.SelectMachine(machines)
	selected_beverage := chooseBeverage(beverages)

	sugarPrice, sugarQuantity := addSugar()
	milkDoses, milkPrice := addMilk(selected_beverage.Name)

	totalPrice := machine.Inventory.CalculatePrice(selected_beverage) + sugarPrice + milkPrice

	fmt.Printf("Total price for %s: CHF %.2f\n", selected_beverage.Name, totalPrice)

	if !machine.Inventory.VerifyStock(selected_beverage) || machine.Inventory.Milk.Quantity < (milkDoses*50) || machine.Inventory.Sugar.Quantity < sugarQuantity {
		return
	}
	processPayment(totalPrice)

	machine.Inventory.UpdateStock(selected_beverage)
	if sugarQuantity > 0 {
		machine.Inventory.Sugar.Quantity -= sugarQuantity
	}
	if milkDoses > 0 {
		machine.Inventory.Milk.Quantity -= milkDoses * 50
	}

	fmt.Printf("Your %s is ready! Enjoy!\n", selected_beverage.Name)
}

func addSugar() (float64, int) {
	fmt.Println("Would you like to add sugar?")
	fmt.Println("1) No sugar")
	fmt.Println("2) Light (5g) - CHF 0.10")
	fmt.Println("3) Medium (10g) - CHF 0.20")
	fmt.Println("4) Heavy (15g) - CHF 0.30")
	choice := validations.GetValidatedInput("> ", []string{"1", "2", "3", "4"})
	switch choice {
	case "2":
		return 0.10, 5
	case "3":
		return 0.20, 10
	case "4":
		return 0.30, 15
	}
	return 0.0, 0
}

func addMilk(beverageName string) (int, float64) {
	if beverageName == "Espresso" {
		return 0, 0.0
	}
	milkChoice := validations.GetValidatedInput("Would you like to add extra milk?\n1) Yes\n2) No", []string{"1", "2"})
	if milkChoice == "2" {
		return 0, 0.0
	}
	doses := validations.GetValidatedNumber("How many doses? (1 to 3)", 1, 3)
	return doses, float64(doses) * 0.05
}

func GenerateTwintCode() string {
	var code strings.Builder
	for i := 0; i < 5; i++ {
		code.WriteByte(chars[rand.Intn(len(chars))])
	}
	return code.String()
}

func processPayment(totalPrice float64) {
	fmt.Printf("Total price: CHF %.2f\n", totalPrice)
	fmt.Println("Please pay using Twint.")
	fmt.Println("Your payment code is:", GenerateTwintCode())
	fmt.Println("(Waiting for payment...)")
	time.Sleep(3 * time.Second)
	fmt.Println("Payment confirmed.")
}

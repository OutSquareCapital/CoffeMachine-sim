package modes_interfaces

import (
	"Nospresso/coffee_machines"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func chooseBeverage() (string, float64, int, int) {
	fmt.Println("Please select your beverage:")
	fmt.Println("1) Espresso - CHF 2.00")
	fmt.Println("2) Cappuccino - CHF 2.50")
	fmt.Println("3) Latte - CHF 2.70 (Small), CHF 3.20 (Medium), CHF 3.70 (Large)")
	beverage := coffee_machines.GetValidatedInput("> ", []string{"1", "2", "3"})
	switch beverage {
	case "1":
		return "Espresso", 2.00, 8, 0
	case "2":
		return "Cappuccino", 2.50, 6, 100
	case "3":
		size := coffee_machines.GetValidatedInput("Select size: 1) Small, 2) Medium, 3) Large", []string{"1", "2", "3"})
		switch size {
		case "1":
			return "Latte (Small)", 2.70, 6, 120
		case "2":
			return "Latte (Medium)", 3.20, 8, 150
		case "3":
			return "Latte (Large)", 3.70, 12, 200
		}
	}
	return "", 0.0, 0, 0
}

func ServeClient(machines []coffee_machines.Machine) {
	if len(machines) == 0 {
		fmt.Println("No machines available.")
		return
	}
	machine := coffee_machines.SelectMachine(machines)
	beverageName, coffeePrice, coffeeNeeded, milkNeeded := chooseBeverage()
	sugarPrice, sugarAmount := addSugar()
	milkDose, milkPrice := addMilk(beverageName)
	totalPrice := coffeePrice + sugarPrice + milkPrice
	if !machine.Inventory.VerifyStock(coffeeNeeded, sugarAmount, milkNeeded+milkDose*50) {
		return
	}
	processPayment(totalPrice)
	machine.Inventory.UpdateStock(coffeeNeeded, sugarAmount, milkNeeded+milkDose*50)
	prepareBeverage(beverageName)
}

func addSugar() (float64, int) {
	fmt.Println("Would you like to add sugar?")
	fmt.Println("1) No sugar")
	fmt.Println("2) Light (5g) - CHF 0.10")
	fmt.Println("3) Medium (10g) - CHF 0.20")
	fmt.Println("4) Heavy (15g) - CHF 0.30")
	choice := coffee_machines.GetValidatedInput("> ", []string{"1", "2", "3", "4"})
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
	milkChoice := coffee_machines.GetValidatedInput("Would you like to add extra milk?\n1) Yes\n2) No", []string{"1", "2"})
	if milkChoice == "2" {
		return 0, 0.0
	}
	doses := coffee_machines.GetValidatedNumber("How many doses? (1 to 3)", 1, 3)
	return doses, float64(doses) * 0.05
}

func GenerateTwintCode() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
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

func prepareBeverage(beverageName string) {
	fmt.Println("Preparing your", beverageName, "...")
	time.Sleep(3 * time.Second)
	fmt.Println("Your", beverageName, "is ready! Enjoy!")
}

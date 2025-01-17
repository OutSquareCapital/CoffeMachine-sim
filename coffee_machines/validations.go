package coffee_machines

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetValidatedQuantity(prompt string) int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(prompt)
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		value, err := strconv.Atoi(input)
		if err == nil && value >= 0 {
			return value
		}
		fmt.Println("Invalid quantity. Try again.")
	}
}

func GetValidatedInput(prompt string, validOptions []string) string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		for _, option := range validOptions {
			if input == option {
				return input
			}
		}
		fmt.Println("Invalid option. Please try again.")
	}
}

func GetValidatedNumber(prompt string, min, max int) int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		value, err := strconv.Atoi(input)
		if err == nil && value >= min && value <= max {
			return value
		}
		fmt.Printf("Invalid input. Please enter a number between %d and %d.\n", min, max)
	}
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func (inv *Inventory) VerifyStock(coffee, sugar, milk int) bool {
	if inv.Coffee.Quantity < coffee {
		fmt.Printf("Error: Insufficient %s to prepare the selected beverage.\n", inv.Coffee.Name)
		return false
	}
	if inv.Sugar.Quantity < sugar {
		fmt.Printf("Error: Insufficient %s to prepare the selected beverage.\n", inv.Sugar.Name)
		return false
	}
	if inv.Milk.Quantity < milk {
		fmt.Printf("Error: Insufficient %s to prepare the selected beverage.\n", inv.Milk.Name)
		return false
	}
	return true
}


package main

  

import (

    "bufio"

    "encoding/csv"

    "errors"

    "fmt"

    "math/rand"

    "os"

    "strconv"

    "strings"

    "time"

)

  

const (

    filename        = "machines.csv"

    ModeClient      = "1"

    ModeAdmin       = "2"

    ModeExit        = "3"

    AdminRestock    = "1"

    AdminUpdatePIN  = "2"

    AdminReturnMenu = "3"

)

  

type Ingredient struct {

    Name         string

    Quantity     int     // Quantité actuelle (en unité ou grammes)

    PricePerUnit float64 // Prix par unité (si applicable)

}

  

type Inventory struct {

    Coffee Ingredient

    Milk   Ingredient

    Sugar  Ingredient

}

  

func (inv *Inventory) DisplayStock() {

    fmt.Printf("Current stock levels:\n")

    fmt.Printf("- %s: %dg\n", inv.Coffee.Name, inv.Coffee.Quantity)

    fmt.Printf("- %s: %.3fL\n", inv.Milk.Name, float64(inv.Milk.Quantity)/1000.0)

    fmt.Printf("- %s: %dg\n", inv.Sugar.Name, inv.Sugar.Quantity)

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

  

func (inv *Inventory) UpdateStock(coffee, sugar, milk int) {

    inv.Coffee.Quantity -= coffee

    inv.Sugar.Quantity -= sugar

    inv.Milk.Quantity -= milk

}

  

func (m *Machine) AddIngredient(ingredient string, amount int) {

    switch ingredient {

    case "milk":

        m.Inventory.Milk.Quantity += amount

    case "sugar":

        m.Inventory.Sugar.Quantity += amount

    case "coffee":

        m.Inventory.Coffee.Quantity += amount

    }

}

  

type Machine struct {

    ID        int

    Pincode   string

    Inventory Inventory

}

  

func selectMachine(machines []Machine) *Machine {

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

            ID:      id + 1,

            Pincode: line[0],

            Inventory: Inventory{

                Coffee: Ingredient{Name: "Coffee", Quantity: coffee},

                Milk:   Ingredient{Name: "Milk", Quantity: milk},

                Sugar:  Ingredient{Name: "Sugar", Quantity: sugar},

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

  

func GenerateTwintCode() string {

    const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

    var code strings.Builder

    for i := 0; i < 5; i++ {

        code.WriteByte(chars[rand.Intn(len(chars))])

    }

    return code.String()

}

  

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

  

func chooseBeverage() (string, float64, int, int) {

    fmt.Println("Please select your beverage:")

    fmt.Println("1) Espresso - CHF 2.00")

    fmt.Println("2) Cappuccino - CHF 2.50")

    fmt.Println("3) Latte - CHF 2.70 (Small), CHF 3.20 (Medium), CHF 3.70 (Large)")

    beverage := GetValidatedInput("> ", []string{"1", "2", "3"})

  

    switch beverage {

    case "1":

        return "Espresso", 2.00, 8, 0

    case "2":

        return "Cappuccino", 2.50, 6, 100

    case "3":

        size := GetValidatedInput("Select size: 1) Small, 2) Medium, 3) Large", []string{"1", "2", "3"})

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

  

func serveClient(machines []Machine) {

    if len(machines) == 0 {

        fmt.Println("No machines available.")

        return

    }

  

    machine := selectMachine(machines)

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

    choice := GetValidatedInput("> ", []string{"1", "2", "3", "4"})

  

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

  

    milkChoice := GetValidatedInput("Would you like to add extra milk?\n1) Yes\n2) No", []string{"1", "2"})

    if milkChoice == "2" {

        return 0, 0.0

    }

  

    doses := GetValidatedNumber("How many doses? (1 to 3)", 1, 3)

    return doses, float64(doses) * 0.05

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

  

func configAdmin(machines []Machine) {

    if len(machines) == 0 {

        fmt.Println("No machines available.")

        return

    }

  

    fmt.Println("Select a machine by ID (1 to", len(machines), ")")

    machineID := GetValidatedNumber("> ", 1, len(machines))

    currentMachine := &machines[machineID-1]

  

    const attempts = 3

    if !ValidatePin(currentMachine, attempts) {

        fmt.Println("Too many failed attempts. Exiting admin mode.")

        return

    }

  

    adminChoice := GetValidatedInput("Access granted.\n1) Restock\n2) Update PIN\n3) Return", []string{"1", "2", "3"})

    switch adminChoice {

    case AdminRestock:

        handleRestock(currentMachine)

    case AdminUpdatePIN:

        UpdatePin(currentMachine)

    case AdminReturnMenu:

        fmt.Println("Returning to main menu.")

    }

}

  

func handleRestock(machine *Machine) {

    machine.Inventory.DisplayStock()

  

    coffee := GetValidatedQuantity("Enter coffee powder (g) to add > ")

    sugar := GetValidatedQuantity("Enter sugar (g) to add > ")

    milk := GetValidatedQuantity("Enter milk (L) to add > ")

  

    machine.Inventory.Coffee.Quantity += coffee

    machine.Inventory.Sugar.Quantity += sugar

    machine.Inventory.Milk.Quantity += milk * 1000

    fmt.Println("Stocks updated successfully.")

}

  

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

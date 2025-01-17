package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const csv_file = "machines.csv"

func main() {
	myApp := app.NewWithID("com.gsemprog.nospresso")
	mainWindow := myApp.NewWindow("Nospresso Cafe")
	mainWindow.Resize(fyne.NewSize(400, 300))

	machines, err := LoadCSV(csv_file)
	if err != nil {
		dialog.ShowError(err, mainWindow)
		return
	}

	clientButton := widget.NewButton("Mode Client", func() {
		displayClientWindow(myApp, machines)
	})
	adminButton := widget.NewButton("Mode Admin", func() {
		displayAdminWindow(myApp, machines)
	})
	exitButton := widget.NewButton("Quitter", func() {
		if err := SaveCSV(csv_file, machines); err != nil {
			dialog.ShowError(err, mainWindow)
		}
		mainWindow.Close()
	})

	mainWindow.SetContent(container.NewVBox(
		widget.NewLabel("Bienvenue à Nospresso Cafe !"),
		clientButton,
		adminButton,
		exitButton,
	))

	mainWindow.ShowAndRun()
}

func displayClientWindow(app fyne.App, machines []Machine) {
	clientWindow := app.NewWindow("Mode Client")
	clientWindow.Resize(fyne.NewSize(400, 300))

	machineSelect := widget.NewSelect(getMachineIDs(machines), func(selected string) {})
	selectButton := widget.NewButton("Sélectionner Machine", func() {
		selectedID, _ := strconv.Atoi(machineSelect.Selected)
		selectedMachine := &machines[selectedID-1]
		serveClientUI(app, selectedMachine)
		clientWindow.Close()
	})
	backButton := widget.NewButton("Retour", func() {
		clientWindow.Close()
	})

	clientWindow.SetContent(container.NewVBox(
		widget.NewLabel("Sélectionnez une machine :"),
		machineSelect,
		selectButton,
		backButton,
	))
	clientWindow.Show()
}

func processPaymentUI(app fyne.App, parentWindow fyne.Window, beverage string) {
	dialog.ShowInformation("Paiement", "Paiement confirmé. Préparation en cours...", parentWindow)
	go func() {
		time.Sleep(3 * time.Second)
		dialog.ShowInformation("Prêt", fmt.Sprintf("Votre %s est prêt. Merci !", beverage), parentWindow)
	}()
}
func serveClientUI(app fyne.App, machine *Machine) {
	clientWindow := app.NewWindow("Mode Client")
	clientWindow.Resize(fyne.NewSize(400, 300))

	container := container.NewVBox() // Conteneur principal
	displayBeverageChoice(app, clientWindow, container, machine)
	clientWindow.SetContent(container)
	clientWindow.Show()
}

func displayBeverageChoice(app fyne.App, window fyne.Window, container *fyne.Container, machine *Machine) {
	beverageSelect := widget.NewSelect([]string{"Espresso", "Cappuccino", "Latte (Small)", "Latte (Medium)", "Latte (Large)"}, nil)
	nextButton := widget.NewButton("Suivant", func() {
		if beverageSelect.Selected == "" {
			dialog.ShowError(errors.New("Veuillez choisir une boisson."), window)
			return
		}
		displaySugarChoice(app, window, container, machine, beverageSelect.Selected)
	})

	container.Objects = []fyne.CanvasObject{
		widget.NewLabel("Étape 1 : Choisissez votre boisson"),
		beverageSelect,
		nextButton,
	}
	container.Refresh()
}

func displaySugarChoice(app fyne.App, window fyne.Window, container *fyne.Container, machine *Machine, beverage string) {
	sugarSelect := widget.NewSelect([]string{"No sugar", "Light (5g)", "Medium (10g)", "Heavy (15g)"}, nil)
	nextButton := widget.NewButton("Suivant", func() {
		if sugarSelect.Selected == "" {
			dialog.ShowError(errors.New("Veuillez choisir un niveau de sucre."), window)
			return
		}
		displayMilkChoice(app, window, container, machine, beverage, sugarSelect.Selected)
	})

	container.Objects = []fyne.CanvasObject{
		widget.NewLabel("Étape 2 : Choisissez votre niveau de sucre"),
		sugarSelect,
		nextButton,
	}
	container.Refresh()
}

func displayMilkChoice(app fyne.App, window fyne.Window, container *fyne.Container, machine *Machine, beverage, sugar string) {
	milkSelect := widget.NewSelect([]string{"No extra milk", "1 dose", "2 doses", "3 doses"}, nil)
	nextButton := widget.NewButton("Suivant", func() {
		if milkSelect.Selected == "" {
			dialog.ShowError(errors.New("Veuillez choisir une option pour le lait."), window)
			return
		}
		displayPayment(app, window, container, machine, beverage, sugar, milkSelect.Selected)
	})

	container.Objects = []fyne.CanvasObject{
		widget.NewLabel("Étape 3 : Choisissez votre dose de lait"),
		milkSelect,
		nextButton,
	}
	container.Refresh()
}

func displayPayment(app fyne.App, window fyne.Window, container *fyne.Container, machine *Machine, beverageName, sugar, milk string) {
	beverage, _, err := getBeverageRequirements(beverageName, sugar, milk)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	if err := machine.Inventory.VerifyStock(beverage.Coffee, beverage.Sugar, beverage.Milk); err != nil {
		dialog.ShowError(err, window)
		return
	}

	totalPrice, err := calculatePrice(beverage, sugar, milk)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	twintCode := GenerateTwintCode()

	confirmButton := widget.NewButton("Payer", func() {
		machine.Inventory.UpdateStock(beverage.Coffee, beverage.Sugar, beverage.Milk)
		dialog.ShowInformation("Succès", fmt.Sprintf("Votre %s est prêt ! Merci pour votre commande.", beverage.Name), window)
		window.Close() // Retour automatique au menu principal
	})

	container.Objects = []fyne.CanvasObject{
		widget.NewLabel(fmt.Sprintf("Étape 4 : Paiement\nTotal : CHF %.2f", totalPrice)),
		widget.NewLabel(fmt.Sprintf("Code Twint : %s", twintCode)),
		confirmButton,
	}
	container.Refresh()
}

func displayAdminWindow(app fyne.App, machines []Machine) {
	adminWindow := app.NewWindow("Mode Admin")
	adminWindow.Resize(fyne.NewSize(400, 300))

	machineSelect := widget.NewSelect(getMachineIDs(machines), func(selected string) {})
	pinEntry := widget.NewEntry()
	pinEntry.SetPlaceHolder("Entrer PIN")

	attemptsRemaining := 3

	loginButton := widget.NewButton("Valider", func() {
		selectedID, _ := strconv.Atoi(machineSelect.Selected)
		selectedMachine := &machines[selectedID-1]

		if ValidatePin(selectedMachine, pinEntry.Text) {
			adminOptionsUI(app, selectedMachine)
			adminWindow.Close()
		} else {
			attemptsRemaining--
			if attemptsRemaining > 0 {
				dialog.ShowError(
					fmt.Errorf("PIN incorrect. %d tentatives restantes.", attemptsRemaining),
					adminWindow,
				)
			} else {
				dialog.ShowError(errors.New("Trop de tentatives échouées. Fermeture de l'application."), adminWindow)
				app.Quit()
			}
		}
	})

	backButton := widget.NewButton("Retour", func() {
		adminWindow.Close()
	})

	adminWindow.SetContent(container.NewVBox(
		widget.NewLabel("Sélectionnez une machine :"),
		machineSelect,
		pinEntry,
		loginButton,
		backButton,
	))
	adminWindow.Show()
}

func adminOptionsUI(app fyne.App, machine *Machine) {
	adminOptionsWindow := app.NewWindow("Options Admin")
	adminOptionsWindow.Resize(fyne.NewSize(400, 300))

	restockButton := widget.NewButton("Restocker Machine", func() {
		handleRestockUI(app, machine)
		adminOptionsWindow.Close()
	})
	updatePinButton := widget.NewButton("Changer PIN", func() {
		handlePinUpdateUI(app, machine)
		adminOptionsWindow.Close()
	})
	backButton := widget.NewButton("Retour", func() {
		adminOptionsWindow.Close()
	})

	adminOptionsWindow.SetContent(container.NewVBox(
		widget.NewLabel("Options administrateur :"),
		restockButton,
		updatePinButton,
		backButton,
	))
	adminOptionsWindow.Show()
}

func handleRestockUI(app fyne.App, machine *Machine) {
	restockWindow := app.NewWindow("Restocker Machine")
	restockWindow.Resize(fyne.NewSize(400, 300))

	// Afficher les stocks actuels
	currentStockLabel := widget.NewLabel(fmt.Sprintf(
		"Stocks actuels :\nCafé : %dg\nSucre : %dg\nLait : %.2fL",
		machine.Inventory.Coffee.Quantity,
		machine.Inventory.Sugar.Quantity,
		float64(machine.Inventory.Milk.Quantity)/1000,
	))

	coffeeEntry := widget.NewEntry()
	coffeeEntry.SetPlaceHolder("Café (g)")
	sugarEntry := widget.NewEntry()
	sugarEntry.SetPlaceHolder("Sucre (g)")
	milkEntry := widget.NewEntry()
	milkEntry.SetPlaceHolder("Lait (L)")

	submitButton := widget.NewButton("Valider", func() {
		coffee, _ := strconv.Atoi(coffeeEntry.Text)
		sugar, _ := strconv.Atoi(sugarEntry.Text)
		milk, _ := strconv.Atoi(milkEntry.Text)

		machine.Inventory.Coffee.Quantity += coffee
		machine.Inventory.Sugar.Quantity += sugar
		machine.Inventory.Milk.Quantity += milk * 1000

		dialog.ShowInformation(
			"Succès",
			fmt.Sprintf("Stocks mis à jour avec succès !\nNouveaux stocks :\nCafé : %dg\nSucre : %dg\nLait : %.2fL",
				machine.Inventory.Coffee.Quantity,
				machine.Inventory.Sugar.Quantity,
				float64(machine.Inventory.Milk.Quantity)/1000,
			),
			restockWindow,
		)

		// Retour automatique au menu principal après mise à jour
		restockWindow.Close()
	})

	backButton := widget.NewButton("Retour", func() {
		restockWindow.Close()
	})

	restockWindow.SetContent(container.NewVBox(
		currentStockLabel,
		coffeeEntry,
		sugarEntry,
		milkEntry,
		submitButton,
		backButton,
	))
	restockWindow.Show()
}

func handlePinUpdateUI(app fyne.App, machine *Machine) {
	pinWindow := app.NewWindow("Changer PIN")
	pinWindow.Resize(fyne.NewSize(400, 200))

	newPinEntry := widget.NewEntry()
	newPinEntry.SetPlaceHolder("Nouveau PIN (6 chiffres)")

	submitButton := widget.NewButton("Valider", func() {
		newPin := newPinEntry.Text
		if err := UpdatePin(machine, newPin); err != nil {
			dialog.ShowError(err, pinWindow)
		} else {
			dialog.ShowInformation("Succès", "PIN mis à jour avec succès", pinWindow)
			pinWindow.Close() // Fermer la fenêtre après succès
		}
	})
	backButton := widget.NewButton("Retour", func() {
		pinWindow.Close()
	})

	pinWindow.SetContent(container.NewVBox(
		widget.NewLabel("Entrer le nouveau PIN :"),
		newPinEntry,
		submitButton,
		backButton,
	))
	pinWindow.Show()
}

package main

import (
	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var retry = make(chan bool)

func main() {
	// Initialize Fyne application
	myApp := app.New()
	mainWindow := myApp.NewWindow("WiZ Bulb Light Controller")

	// Create initial UI components
	statusLabel := widget.NewLabel("Discovering WiZ devices...")
	retryButton := widget.NewButton("Retry", func() {
		log.Println("Retry button clicked")
		retry <- true
	})

	// Add components to the main window
	mainWindow.SetContent(container.NewVBox(
		statusLabel,
		retryButton,
	))

	// Show the main window
	mainWindow.Show()

	// Attempt to discover WiZ devices
	go func() {
		devices := discoverWiZDevices()

		for len(devices) == 0 {
			statusLabel.SetText("No WiZ devices found.")
			retryButton.Show()

			// Wait for retry button to be clicked
			<-retry

			// Retry discovery
			devices = discoverWiZDevices()
		}

		// Initialize WizFleet
		fleet := NewWizFleet()
		fleet.Devices = devices

		// Start monitoring devices
		fleet.Start()

		// Create and show the main window with actual content
		mainWindow := createMainWindow(myApp, fleet)
		mainWindow.CenterOnScreen()
		mainWindow.Show()
	}()

	// Run the application
	myApp.Run()
}

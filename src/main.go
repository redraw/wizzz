package main

import (
	"log"

	"fyne.io/fyne/v2/app"
)

func main() {
	// Discover WiZ devices
	devices := discoverWiZDevices()
	if len(devices) == 0 {
		log.Fatal("No WiZ devices found on the network")
	}

	// Initialize WizFleet
	fleet := NewWizFleet()
	fleet.Devices = devices

	// Start monitoring devices
	fleet.Start()

	// Initialize and run the GUI
	myApp := app.New()
	mainWindow := createMainWindow(myApp, fleet)
	mainWindow.Show()

	// Run the application
	myApp.Run()

	// Close UDP connections
	for _, device := range devices {
		device.conn.Close()
	}
}

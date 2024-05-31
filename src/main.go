package main

import (
	"flag"
	"time"

	"fyne.io/fyne/v2/app"
	log "github.com/sirupsen/logrus"
)

const title = "WiZ LED controller"

func main() {
	logLevel := flag.String("loglevel", "info", "set the log level (debug, info, warn, error, fatal, panic)")
	discoveryTimeout := flag.Int("discovery-timeout", 2, "set the discovery timeout in seconds")
	flag.Parse()

	if level, err := log.ParseLevel(*logLevel); err == nil {
		log.SetLevel(level)
	}

	// Initialize Fyne application
	myApp := app.New()
	mainWindow := myApp.NewWindow(title)
	fleet := NewWizFleet()

	deviceCh := make(chan *WizDevice)
	go waitForDevices(mainWindow, time.Second*time.Duration(*discoveryTimeout), deviceCh)

	go func() {
		for device := range deviceCh {
			fleet.AddDevice(device)
		}
		start(mainWindow, fleet)
	}()

	// Run the application
	myApp.Run()
}

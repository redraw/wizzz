package main

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/lusingander/colorpicker"
	log "github.com/sirupsen/logrus"
)

func waitForDevices(mainWindow fyne.Window, timeout time.Duration, ch chan *WizDevice) {
	retry := make(chan bool)

	statusLabel := widget.NewLabel("Discovering WiZ devices...")
	retryButton := widget.NewButton("Retry", func() {
		retry <- true
	})

	mainWindow.SetContent(container.NewVBox(
		statusLabel,
		retryButton,
	))

	mainWindow.Show()

	// Attempt to discover WiZ devices
	devices := discoverWiZDevices(timeout)

	for len(devices) == 0 {
		statusLabel.SetText("No WiZ devices found.")
		retryButton.Show()

		// Wait for retry button to be clicked
		<-retry

		// Retry discovery
		devices = discoverWiZDevices(timeout)
	}

	// Send discovered devices to the main function
	for _, device := range devices {
		ch <- device
	}
	close(ch)

	mainWindow.Hide()
}

func start(mainWindow fyne.Window, fleet *WizFleet) {
	fleet.Start()

	// Controls
	switchButton := widget.NewCheck("Power", func(on bool) {
		go fleet.SetPower(on)
	})

	brightnessSlider := widget.NewSlider(10, 100)
	brightnessSlider.SetValue(50)
	brightnessSlider.OnChanged = func(value float64) {
		go fleet.SetBrightness(value)
	}

	temperatureSlider := widget.NewSlider(2200, 6500)
	temperatureSlider.SetValue(4000)
	temperatureSlider.OnChanged = func(value float64) {
		go fleet.SetTemperature(value)
	}

	colorPicker := colorpicker.New(200, colorpicker.StyleHue)
	colorPicker.SetOnChanged(func(c color.Color) {
		go fleet.SetColor(c)
	})

	deviceSelector := widget.NewSelect(nil, func(value string) {
		device := fleet.Select(value)
		log.Debugf("Selected device: %+v", device)
		if device != nil {
			if state, err := device.GetState(); err == nil {
				switchButton.SetChecked(state.State)
				brightnessSlider.SetValue(state.Dimming)
				temperatureSlider.SetValue(state.Temp)
			}
		}
	})

	options := make([]string, 0, len(fleet.Devices))
	options = append(options, "All")
	for _, device := range fleet.Devices {
		options = append(options, device.IP)
	}
	deviceSelector.SetOptions(options)
	deviceSelector.SetSelectedIndex(1)

	rssiLabel := widget.NewLabel("RSSI: -")

	// Container
	content := container.NewVBox(
		deviceSelector,
		rssiLabel,
		switchButton,
		widget.NewLabel("Brightness:"),
		brightnessSlider,
		widget.NewLabel("Temperature:"),
		temperatureSlider,
		widget.NewLabel("Color:"),
		container.NewVBox(colorPicker),
	)

	mainWindow.SetContent(content)
	mainWindow.CenterOnScreen()
	mainWindow.Show()

	// Update RSSI every 1 second
	go func() {
		for {
			if fleet.SelectedDevice != nil {
				rssiLabel.SetText(fmt.Sprintf("RSSI: %v", fleet.SelectedDevice.State.Rssi))
			}
			time.Sleep(time.Second)
		}
	}()
}

package main

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/lusingander/colorpicker"
)

func createMainWindow(myApp fyne.App, fleet *WizFleet) fyne.Window {
	mainWindow := myApp.NewWindow("WiZ Light Bulb Controller")

	switchButton := widget.NewCheck("Power", func(on bool) {
		go fleet.SetPower(on)
	})

	brightnessSlider := widget.NewSlider(0, 100)
	brightnessSlider.OnChanged = func(value float64) {
		go fleet.SetBrightness(int(value))
	}

	temperatureSlider := widget.NewSlider(2200, 6500)
	temperatureSlider.OnChanged = func(value float64) {
		go fleet.SetTemperature(int(value))
	}

	colorPicker := colorpicker.New(200, colorpicker.StyleHue)
	colorPicker.SetOnChanged(func(c color.Color) {
		go fleet.SetColor(c)
	})

	rssiLabel := widget.NewLabel("RSSI: -")

	deviceSelector := widget.NewSelect(nil, func(value string) {
		fleet.Select(value)
		if fleet.SelectedDevice != nil {
			if state := fleet.SelectedDevice.State; state != nil {
				switchButton.SetChecked(state["state"].(bool))
				brightnessSlider.SetValue(float64(state["dimming"].(float64)))
				temperatureSlider.SetValue(float64(state["temp"].(float64)))
			}
		}
	})

	options := make([]string, 0, len(fleet.Devices))
	options = append(options, "All")
	for _, device := range fleet.Devices {
		options = append(options, device.IP)
	}
	deviceSelector.SetOptions(options)
	deviceSelector.SetSelected(options[1])

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

	go func() {
		for {
			if fleet.SelectedDevice != nil {
				rssiLabel.SetText(fmt.Sprintf("RSSI: %d", int(fleet.SelectedDevice.State["rssi"].(float64))))
			}
			time.Sleep(1 * time.Second)
		}
	}()

	return mainWindow
}

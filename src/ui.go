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

func start(mainWindow fyne.Window, deviceCh chan *WizDevice) {
	fleet := NewWizFleet()
	fleet.Start()

	mainWindow.SetContent(container.NewVBox(
		widget.NewLabel("Discovering WiZ devices..."),
	))

	mainWindow.Show()

	for len(fleet.Devices) == 0 {
		device := <-deviceCh
		if device != nil {
			fleet.AddDevice(device)
		}
	}

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
				switchButton.SetChecked(true)
				brightnessSlider.SetValue(state.Dimming)
				temperatureSlider.SetValue(state.Temp)
			}
		}
	})

	options := make([]string, 0)
	options = append(options, "All")
	options = append(options, fleet.SelectedDevice.IP)
	deviceSelector.SetOptions(options)
	deviceSelector.SetSelectedIndex(1)

	// Update device list until deviceCh is consumed
	go func() {
		for device := range deviceCh {
			fleet.AddDevice(device)
			options = append(options, device.IP)
			deviceSelector.SetOptions(options)
		}
	}()

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

package main

import (
	"image/color"
	"log"
	"time"
)

type WizFleet struct {
	Devices        []*WizDevice
	SelectedDevice *WizDevice
}

func NewWizFleet() *WizFleet {
	return &WizFleet{
		Devices: make([]*WizDevice, 0),
	}
}

func (wf *WizFleet) Select(ip string) *WizDevice {
	device := getSelectedDevice(ip, wf.Devices)
	wf.SelectedDevice = device
	return device
}

func (wf *WizFleet) Start() {
	for _, device := range wf.Devices {
		go wf.monitorDevice(device)
	}
}

func (wf *WizFleet) monitorDevice(device *WizDevice) {
	for {
		err := device.sendCommand("getPilot", nil)
		if err != nil {
			log.Printf("Error getting state for device %s: %v", device.IP, err)
		}
		time.Sleep(1 * time.Second)
	}
}

func (wf *WizFleet) SetPower(on bool) error {
	if wf.SelectedDevice != nil {
		return wf.SelectedDevice.SetPower(on)
	}
	for _, device := range wf.Devices {
		err := device.SetPower(on)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wf *WizFleet) SetBrightness(value int) error {
	if wf.SelectedDevice != nil {
		return wf.SelectedDevice.SetBrightness(value)
	}
	for _, device := range wf.Devices {
		err := device.SetBrightness(value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wf *WizFleet) SetTemperature(value int) error {
	if wf.SelectedDevice != nil {
		return wf.SelectedDevice.SetTemperature(value)
	}
	for _, device := range wf.Devices {
		err := device.SetTemperature(value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wf *WizFleet) SetColor(rgb color.Color) error {
	if wf.SelectedDevice != nil {
		return wf.SelectedDevice.SetColor(rgb)
	}
	for _, device := range wf.Devices {
		err := device.SetColor(rgb)
		if err != nil {
			return err
		}
	}
	return nil
}

func getSelectedDevice(selected string, devices []*WizDevice) *WizDevice {
	if selected == "" || selected == "All" {
		return nil
	}

	for _, device := range devices {
		if device.IP == selected {
			return device
		}
	}

	return nil
}

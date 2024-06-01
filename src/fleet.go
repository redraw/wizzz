package main

import (
	"image/color"
	"time"

	log "github.com/sirupsen/logrus"
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

func (wf *WizFleet) AddDevice(device *WizDevice) {
	wf.Devices = append(wf.Devices, device)
}

func (wf *WizFleet) Select(ip string) *WizDevice {
	device := getSelectedDevice(ip, wf.Devices)
	wf.SelectedDevice = device
	return device
}

func (wf *WizFleet) Start() {
	log.Debug("Start monitoring...")
	for _, device := range wf.Devices {
		go wf.monitorDevice(device)
	}
}

func (wf *WizFleet) monitorDevice(device *WizDevice) {
	for {
		if state, err := device.GetState(); err == nil {
			log.Debugf("Device %s state: %+v", device.IP, state)
		}
		time.Sleep(5 * time.Second)
	}
}

func (wf *WizFleet) SetPower(on bool) error {
	if wf.SelectedDevice != nil {
		if _, err := wf.SelectedDevice.SetPower(on); err != nil {
			return err
		}
	} else {
		for _, device := range wf.Devices {
			go device.SetPower(on)
		}
	}
	return nil
}

func (wf *WizFleet) SetBrightness(value float64) error {
	if wf.SelectedDevice != nil {
		if _, err := wf.SelectedDevice.SetBrightness(value); err != nil {
			return err
		}
	} else {
		for _, device := range wf.Devices {
			go device.SetBrightness(value)
		}
	}
	return nil
}

func (wf *WizFleet) SetTemperature(value float64) error {
	if wf.SelectedDevice != nil {
		if _, err := wf.SelectedDevice.SetTemperature(value); err != nil {
			return err
		}
	} else {
		for _, device := range wf.Devices {
			go device.SetTemperature(value)
		}
	}
	return nil
}

func (wf *WizFleet) SetColor(rgb color.Color) error {
	if wf.SelectedDevice != nil {
		if _, err := wf.SelectedDevice.SetColor(rgb); err != nil {
			return err
		}
	} else {
		for _, device := range wf.Devices {
			go device.SetColor(rgb)
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

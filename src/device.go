package main

import (
	"encoding/json"
	"image/color"
	"log"
	"net"
)

type WizDevice struct {
	IP    string
	MAC   string
	State map[string]interface{}
	conn  *net.UDPConn
}

func NewWizDevice(ip, mac string) (*WizDevice, error) {
	addr, err := net.ResolveUDPAddr("udp", ip+":38899")
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	device := &WizDevice{
		IP:    ip,
		MAC:   mac,
		conn:  conn,
		State: make(map[string]interface{}),
	}

	return device, nil
}

func (wd *WizDevice) sendCommand(method string, params map[string]interface{}) error {
	log.Println(">", method, params)

	payload := map[string]interface{}{
		"method": method,
		"params": params,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = wd.conn.Write(body)
	if err != nil {
		return err
	}

	response := make([]byte, 4096)
	n, err := wd.conn.Read(response)
	if err != nil {
		return err
	}

	var respData struct {
		Result map[string]interface{} `json:"result"`
	}
	err = json.Unmarshal(response[:n], &respData)
	if err != nil {
		return err
	}

	log.Println("<", respData)

	if method == "getPilot" {
		wd.State = respData.Result
	}

	return nil
}

func (wd *WizDevice) SetPower(on bool) error {
	state := map[string]interface{}{"state": on}
	return wd.sendCommand("setPilot", state)
}

func (wd *WizDevice) SetBrightness(value int) error {
	state := map[string]interface{}{"dimming": value}
	return wd.sendCommand("setPilot", state)
}

func (wd *WizDevice) SetTemperature(value int) error {
	state := map[string]interface{}{"temp": value}
	return wd.sendCommand("setPilot", state)
}

func (wd *WizDevice) SetColor(rgb color.Color) error {
	r, g, b, _ := rgb.RGBA()
	state := map[string]interface{}{"r": int(r >> 8), "g": int(g >> 8), "b": int(b >> 8)}
	return wd.sendCommand("setPilot", state)
}

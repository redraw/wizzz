package main

import (
	"encoding/json"
	"errors"
	"image/color"
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type WizDevice struct {
	IP    string
	MAC   string
	State WizState
	conn  *net.UDPConn
	mu    sync.Mutex
}

type WizState struct {
	Mac     string  `json:"mac"`
	Dimming float64 `json:"dimming"`
	State   bool    `json:"state"`
	Rssi    float64 `json:"rssi"`
	Temp    float64 `json:"temp"`
	R       int     `json:"r"`
	G       int     `json:"g"`
	B       int     `json:"b"`
}

type WizParams map[string]interface{}

type WizPayload struct {
	Method string    `json:"method"`
	Params WizParams `json:"params"`
}

func (wd *WizDevice) Close() {
	wd.conn.Close()
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
		IP:   ip,
		MAC:  mac,
		conn: conn,
	}

	return device, nil
}

func (wd *WizDevice) sendCommand(payload WizPayload) (interface{}, error) {
	wd.mu.Lock()
	defer wd.mu.Unlock()

	log.Debugf("> %+v", payload)

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	_, err = wd.conn.Write(body)
	if err != nil {
		return nil, err
	}

	bytes := make([]byte, 4096)
	if err := wd.conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		return nil, err
	}

	n, err := wd.conn.Read(bytes)
	if err != nil {
		return nil, err
	}

	var response struct {
		Method string `json:"method"`
	}
	if err := json.Unmarshal(bytes[:n], &response); err != nil {
		return nil, err
	}

	switch response.Method {
	case "getPilot":
		var data struct {
			Result WizState `json:"result"`
		}
		if err := json.Unmarshal(bytes[:n], &data); err != nil {
			return nil, err
		}
		log.Debugf("< %+v", data.Result)
		wd.State = data.Result
		return data.Result, nil

	case "setPilot":
		var data struct {
			Result struct {
				Success bool `json:"success"`
			} `json:"result"`
		}
		if err := json.Unmarshal(bytes[:n], &data); err != nil {
			return nil, err
		}
		log.Debugf("< %+v", data.Result)
		return data.Result.Success, nil
	}

	return nil, errors.New("unknown response")
}

func (wd *WizDevice) GetState() (WizState, error) {
	request := WizPayload{Method: "getPilot"}

	res, err := wd.sendCommand(request)
	if err != nil {
		return WizState{}, err
	}

	return res.(WizState), nil
}

func (wd *WizDevice) SetPower(on bool) (bool, error) {
	request := WizPayload{
		Method: "setPilot",
		Params: WizParams{"state": on},
	}

	res, err := wd.sendCommand(request)
	if err != nil {
		return false, err
	}

	return res.(bool), nil
}

func (wd *WizDevice) SetBrightness(value float64) (bool, error) {
	request := WizPayload{
		Method: "setPilot",
		Params: WizParams{"dimming": value},
	}

	res, err := wd.sendCommand(request)
	if err != nil {
		return false, err
	}

	return res.(bool), nil
}

func (wd *WizDevice) SetTemperature(value float64) (bool, error) {
	request := WizPayload{
		Method: "setPilot",
		Params: WizParams{"temp": value},
	}

	res, err := wd.sendCommand(request)
	if err != nil {
		return false, err
	}

	return res.(bool), nil
}

func (wd *WizDevice) SetColor(rgb color.Color) (bool, error) {
	r, g, b, _ := rgb.RGBA()

	request := WizPayload{
		Method: "setPilot",
		Params: WizParams{
			"r": int(r >> 8),
			"g": int(g >> 8),
			"b": int(b >> 8),
		},
	}

	res, err := wd.sendCommand(request)
	if err != nil {
		return false, err
	}

	return res.(bool), nil
}

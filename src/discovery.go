package main

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

const (
	discoveryPort = 38899
	discoveryMsg  = `{"method":"registration","params":{"phoneMac":"AAAAAAAAAAAA","phoneIp":"1.2.3.4","register":false,"id":1}}`
)

func discoverWiZDevices() []*WizDevice {
	var devices []*WizDevice

	localAddr := &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 0,
	}

	serverAddr := net.UDPAddr{
		IP:   net.ParseIP("255.255.255.255"),
		Port: discoveryPort,
	}

	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Fatalf("Error setting up discovery: %v", err)
	}
	defer conn.Close()

	timeout := time.After(1 * time.Second)

	for {
		select {
		case <-timeout:
			log.Println("Discovery timeout reached")
			return devices
		default:
			log.Println("Sending discovery packet...")
			_, err := conn.WriteToUDP([]byte(discoveryMsg), &serverAddr)
			if err != nil {
				log.Printf("Error sending discovery packet: %v", err)
				continue
			}

			buf := make([]byte, 4096)
			conn.SetReadDeadline(time.Now().Add(2 * time.Second))
			n, addr, err := conn.ReadFromUDP(buf)
			if err != nil {
				continue
			}

			var data map[string]interface{}
			err = json.Unmarshal(buf[:n], &data)
			if err != nil {
				continue
			}

			if result, ok := data["result"].(map[string]interface{}); ok {
				if success, ok := result["success"].(bool); ok && success {
					mac := result["mac"].(string)
					ip := addr.IP.String()
					if !isDeviceAlreadyDiscovered(devices, mac) {
						device, err := NewWizDevice(ip, mac)
						if err != nil {
							continue
						}
						devices = append(devices, device)
					}
				}
			}
		}
	}
}

func isDeviceAlreadyDiscovered(devices []*WizDevice, mac string) bool {
	for _, device := range devices {
		if device.MAC == mac {
			return true
		}
	}
	return false
}
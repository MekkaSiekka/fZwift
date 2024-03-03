/*
https://martin-ueding.de/posts/heart-rate-monitor-with-python/
*/
package main

import (
	"fmt"
	"time"

	"github.com/MekkaSiekka/fZwift/adapters"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func main() {
	// Enable BLE interface.
	must("enable BLE stack", adapter.Enable())

	deviceScanner := adapters.NewDeviceScanner()
	time.Sleep(time.Second * 15)
	devices := deviceScanner.GetAllDevices()
	for key, val := range devices {
		fmt.Printf("Key %v, address %v, name %v\n", key, val.Address, val.Name)
	}
	// Start scanning.
	println("scanning...")
	// err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
	// 	println("found device:", device.Address.String(), device.RSSI, device.LocalName())
	// })
	// 	must("start scan", err)
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}

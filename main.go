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
	time.Sleep(time.Second * 5)
	devices := deviceScanner.GetAllDevices()
	for key, val := range devices {
		fmt.Printf("Key %v, address %v, name %v\n", key, val.Address, val.Name)
		var device *bluetooth.Device
		fmt.Printf("Trying to connect %v\n", val.Address)
		device, err := adapter.Connect(val.Address, bluetooth.ConnectionParams{})
		if err != nil {
			println(err.Error())
			continue
			//return fmt.Errorf("connect device : %v", err)
		}
		fmt.Printf("connected to %v\n", val.Address.String())
		//controlChar := "00002a66-0000-1000-8000-00805f9b34fb"
		services, err := device.DiscoverServices(
			[]bluetooth.UUID{bluetooth.ServiceUUIDFitnessMachine})
		if err != nil {
			fmt.Printf("discover char for device %v: %v\n", val.Name, err)
			continue
		}
		if len(services) == 0 {
			fmt.Printf("This device does not have service\n")
			continue
		}
		service := services[0]
		chars, err := service.DiscoverCharacteristics([]bluetooth.UUID{bluetooth.CharacteristicUUIDDFUControlPoint})
		if err != nil {
			fmt.Printf("find char %v", err)
		}
		if len(chars) == 0 {
			fmt.Printf("cannot find char")
		}
		fmt.Printf("=================GOOD And Found Controller =======================\n")

	}

}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}

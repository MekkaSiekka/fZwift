package adapters

import (
	"fmt"

	"tinygo.org/x/bluetooth"
)

type HRManagerImpl struct {
	devices map[string]bool // address is this type
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}

var (
	adapter                     = bluetooth.DefaultAdapter
	heartRateServiceUUID        = bluetooth.ServiceUUIDHeartRate
	heartRateCharacteristicUUID = bluetooth.CharacteristicUUIDHeartRateMeasurement
)

func newHRManager() *HRManagerImpl {

	return &HRManagerImpl{}
}

/*
Ultimately just list all devices that are HRM

Should be async
*/
func (hr *HRManagerImpl) addDevicesToList() error {
	println("enabling")

	// // Enable BLE interface.
	//must("enable BLE stack", adapter.Enable())

	ch := make(chan bluetooth.ScanResult, 1)

	// Start scanning.
	println("scanning...")
	err := adapter.Scan(
		func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
			println("found device:", result.Address.String(), result.RSSI, result.LocalName())
			adapter.StopScan()
			ch <- result
		})

	var device *bluetooth.Device

	// While loop to wait for devices
	select {
	case result := <-ch:
		println("Trying to connect ", result.Address.String())
		device, err = adapter.Connect(result.Address, bluetooth.ConnectionParams{})
		if err != nil {
			println(err.Error())

			return fmt.Errorf("connect device : %v", err)
		}
		println("connected to ", result.Address.String())
	}

	// get services
	println("discovering services/characteristics")
	srvcs, err := device.DiscoverServices([]bluetooth.UUID{heartRateServiceUUID})
	// must("discover services", err)

	srvc := srvcs[0]

	println("found service", srvc.UUID().String())

	chars, err := srvc.DiscoverCharacteristics([]bluetooth.UUID{heartRateCharacteristicUUID})
	if err != nil {
		println(err)
	}

	if len(chars) == 0 {
		panic("could not find heart rate characteristic")
	}

	char := chars[0]
	println("found characteristic", char.UUID().String())

	// Setup notification and callback upon receiving data
	char.EnableNotifications(func(buf []byte) {
		println("data:", uint8(buf[1]))
	})

	select {}
}

package adapters

import (
	"fmt"
	"strings"

	"tinygo.org/x/bluetooth"
)

type DeviceInfo struct {
	Address string
	Name    string
}

type DeviceScannerImpl struct {
	devices      map[string]DeviceInfo // address is this type
	discoverChan chan bluetooth.ScanResult
}

func newDeviceScanner() *DeviceScannerImpl {
	ds := &DeviceScannerImpl{
		devices:      map[string]DeviceInfo{},
		discoverChan: make(chan bluetooth.ScanResult, 1),
	}
	go ds.startScan()
	return ds
}

func (ds *DeviceScannerImpl) GetAllDevices() map[string]DeviceInfo {
	return ds.devices
}

// Async function to scan and add to map
func (ds *DeviceScannerImpl) startScan() {
	// Register the callback function so whenever
	// scan has some result, do something
	adapter.Scan(
		ds.handleDiscoveredDevice,
	)
	adapter.StopScan()
}

func (ds *DeviceScannerImpl) handleDiscoveredDevice(
	btAdatper *bluetooth.Adapter,
	btScanResult bluetooth.ScanResult,
) {
	fmt.Printf(
		"found device: %v, %v, %v\n",
		btScanResult.Address.String(),
		btScanResult.RSSI,
		btScanResult.LocalName(),
	)
	di := DeviceInfo{
		Address: btScanResult.Address.String(),
	}
	trimmedName := strings.TrimSpace(btScanResult.LocalName())
	if len(trimmedName) > 0 {
		fmt.Printf("Name is %v\n", trimmedName)
		di.Name = trimmedName
	}
	ds.updateDeviceInfo(di)
}

func (ds *DeviceScannerImpl) updateDeviceInfo(di DeviceInfo) {
	_, exist := ds.devices[di.Address]
	if !exist {
		ds.devices[di.Address] = di
		return
	}
	// Case the device is already discovered
	if len(di.Name) > 0 {
		ds.devices[di.Address] = di
	}

}

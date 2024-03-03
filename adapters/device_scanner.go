package adapters

type DeviceScanner interface {
	GetAllDevices() map[string]DeviceInfo
}

func NewDeviceScanner() DeviceScanner {
	return newDeviceScanner()
}

/*
https://martin-ueding.de/posts/heart-rate-monitor-with-python/
*/
package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/MekkaSiekka/fZwift/adapters"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func getNthBit(byteArray []byte, n int) int {
	byteIndex := n / 8
	bitIndex := 7 - (n % 8)
	if byteIndex >= len(byteArray) {
		return -1 // Out of bounds
	}
	bit := (byteArray[byteIndex] >> uint(bitIndex)) & 1
	return int(bit)
}

func main() {
	// Enable BLE interface.
	must("enable BLE stack", adapter.Enable())

	deviceScanner := adapters.NewDeviceScanner()
	time.Sleep(time.Second * 5)
	devices := deviceScanner.GetAllDevices()
	// buffer to retrieve characteristic data
	buf := make([]byte, 8)
	for key, val := range devices {
		fmt.Printf("Key %v, address %v, name %v\n", key, val.Address, val.Name)
		if val.Name != "Garmin Flux 58884" {
			continue
		}
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
			fmt.Printf("discover char for device fail %v: %v\n", val.Name, err)
			continue
		}
		if len(services) == 0 {
			fmt.Printf("This device does not have service\n")
			continue
		}
		for _, service := range services {
			chars, err := service.DiscoverCharacteristics(
				[]bluetooth.UUID{bluetooth.CharacteristicUUIDFitnessMachineFeature},
			)
			if err != nil {
				fmt.Printf("found error char %v", err)
				continue
			}
			if len(chars) == 0 {
				fmt.Printf("cannot find char\n")
				continue
			}
			char := chars[0]
			fmt.Printf("=================GOOD And Found Controller =======================\n")

			n, err := char.Read(buf)
			if err != nil {
				fmt.Printf("    %v\n", err.Error())
				continue
			}
			fmt.Printf("    data bytes %v\n", strconv.Itoa(n))
			//println("    value =  \n", string(buf[:n]))
			fmt.Printf("%08b\n", buf[0:n])

			for _, b := range buf {
				for i := 0; i < 8; i++ {
					bit := (b >> uint(i)) & 1
					fmt.Print(bit)
				}
				fmt.Println()
			}

			parser := adapters.FitnessMachineChar{}
			parser.ParseCharBuffer(buf)

			fmt.Printf("Finish print \n")
			// Requst control
			chars, err = service.DiscoverCharacteristics(
				[]bluetooth.UUID{bluetooth.CharacteristicUUIDFitnessMachineControlPoint},
			)
			if err != nil {
				fmt.Printf("found error char %v", err)
				continue
			}
			if len(chars) == 0 {
				fmt.Printf("cannot find char\n")
				continue
			}
			char = chars[0]
			char.EnableNotifications(func(buf []byte) {
				for idx, b := range buf {
					fmt.Printf("callback: %v : %v\n", idx, b)
				}
			})
			fmt.Printf("=================GOOD And Found Control Point =======================\n")
			bytesWritten, err := char.Write([]byte{0x00})
			if err != nil {
				fmt.Printf("Cannot write byte %v", err)
				continue
			}
			// fmt.Printf("Bytes writte : %v \n", bytesWritten)
			// bytesWritten, err = char.Write([]byte{0x08})
			// if err != nil {
			// 	fmt.Printf("Cannot write byte get response %v", err)
			// 	continue
			// }
			fmt.Printf("Bytes written : %v \n", bytesWritten)
			time.Sleep(time.Second * 10)
			for {
				bytesWritten, err = char.Write([]byte{0x05, 0x64, 0x00})
				if err != nil {
					fmt.Printf("Cannot write byte get response %v", err)
					continue
				} else {
					fmt.Printf("Bytes writte to contorl 1 : %v \n", bytesWritten)
				}
				time.Sleep(10 * time.Second)
				bytesWritten, err = char.Write([]byte{0x05, 0xC8, 0x00})
				if err != nil {
					fmt.Printf("Cannot write byte get response %v", err)
					continue
				} else {
					fmt.Printf("Bytes writte to control 2: %v \n", bytesWritten)
				}
				time.Sleep(10 * time.Second)
			}

		}
	}

}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}

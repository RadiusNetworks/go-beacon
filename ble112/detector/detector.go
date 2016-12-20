package main

import (
	"fmt"
	"github.com/RadiusNetworks/go-beacon/ble112"
	"os"
)

func main() {
	devices := ble112.Devices()
	if len(devices) == 0 {
		fmt.Printf("No BLE112 devices found!\n")
		os.Exit(1)
	} else {
		for _, device := range devices {
			fmt.Printf("%v => %v\n", device.MacAddress, device.Port)
		}
	}
}

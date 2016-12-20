package main

import (
	"fmt"
	"github.com/RadiusNetworks/go-beacon/ble112"
	"os"
	"sync"
	"time"
)

func main() {
	paths, _ := ble112.DevicePaths()
	devices := make(chan *ble112.Device, 30)
	var wg sync.WaitGroup
	for _, port := range paths {
		wg.Add(1)
		go func(port string) {
			defer wg.Done()
			ch := make(chan *ble112.Device)
			go func() {
				device, _ := ble112.NewDevice(port)
				ch <- device
			}()
			timeout := time.NewTimer(500 * time.Millisecond)
			select {
			case device := <-ch:
				if device != nil {
					devices <- device
				}
				return
			case <-timeout.C:
				return
			}
		}(port)
	}
	wg.Wait()
	close(devices)

	count := 0
	for device := range devices {
		fmt.Printf("%v => %v\n", device.MacAddress, device.Port)
		count += 1
	}
	if count == 0 {
		fmt.Printf("No BLE112 devices found!\n")
		os.Exit(1)
	}
}

package advertiser

import (
	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/linux"
)

func defaultDevice() (ble.Device, error) {
	device, err := linux.NewDevice()
	if err != nil {
		return nil, err
	}
	return device, nil
}

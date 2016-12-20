package advertiser

import (
	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/darwin"
)

func defaultDevice() (ble.Device, error) {
	device, err := darwin.NewDevice()
	if err != nil {
		return nil, err
	}
	return device, nil
}

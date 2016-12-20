package advertiser

import (
	"errors"
	"github.com/currantlabs/ble"
)

func defaultDevice() (ble.Device, error) {
	return nil, errors.New("Advertising not supported on Windows")
}

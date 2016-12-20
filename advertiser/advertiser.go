package advertiser

import (
	"github.com/currantlabs/ble"
	"golang.org/x/net/context"
)

// An Advertisement contains the bytes of a beacon advertisement.
type Advertisement []byte

// An Advertiser represents hardware that can advertise as a beacon.
type Advertiser interface {
	AdvertiseMfgData(id uint16, ad Advertisement)
	AdvertiseServiceData(id uint16, ad Advertisement)
	StopAdvertising()
}

type advertiser struct {
	device ble.Device
	ctx    context.Context
	cancel context.CancelFunc
	done   chan bool
}

// New returns a new Advertiser using the default BLE hardware.
func New() (Advertiser, error) {
	device, err := defaultDevice()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &advertiser{
		device: device,
		ctx:    ctx,
		cancel: cancel,
		done:   make(chan bool),
	}, nil
}

// AdvertiseMfgData advertises manufacturer data with the given mfg id
func (a *advertiser) AdvertiseMfgData(id uint16, ad Advertisement) {
	go func() {
		a.device.AdvertiseMfgData(a.ctx, id, ad[2:])
		a.done <- true
	}()
}

// AdvertiseServiceData advertises service data given a 16bit UUID
func (a *advertiser) AdvertiseServiceData(id uint16, ad Advertisement) {
	go func() {
		a.device.AdvertiseServiceData16(a.ctx, id, ad[2:])
		a.done <- true
	}()
}

// StopAdvertising stops the hardware from advertising as a beacon.
func (a *advertiser) StopAdvertising() {
	a.cancel()
	<-a.done
}

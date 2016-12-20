package ble112

import (
	"bytes"
	"errors"
	"github.com/RadiusNetworks/go-beacon"
	"github.com/RadiusNetworks/go-beacon/advertiser"
	"github.com/tarm/serial"
	"path/filepath"
	"runtime"
	"time"
)

// Device represents a USB connected BLE112 which can be used for
// BLE scanning or advertising.
type Device struct {
	Port       string
	MacAddress *beacon.MacAddress
	f          *serial.Port
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const (
	BG_COMMAND              = byte(0)
	BG_MSG_CLASS_SYSTEM     = byte(0)
	BG_MSG_CLASS_CONNECTION = byte(3)
	BG_MSG_CLASS_GAP        = byte(6)
	BG_GET_ADDRESS          = byte(2)
	BG_DISCONNECT           = byte(0)
	BG_SET_MODE             = byte(1)
	BG_DISCOVER             = byte(2)
	BG_DISCOVER_STOP        = byte(4)
	BG_SCAN_PARAMS          = byte(7)
	BG_GAP_NON_DISCOVERABLE = byte(0)
	BG_GAP_NON_CONNECTABLE  = byte(0)
	BG_GAP_CONNECTABLE      = byte(2)
	BG_GAP_DISCOVER_ALL     = byte(2)
	BG_GAP_USER_DATA        = byte(4)
	BG_GAP_SET_ADV_PARAM    = byte(8)
	BG_GAP_SET_ADV_DATA     = byte(9)
	BG_EVENT                = byte(0x80)
)

var NULL_DATA = make([]byte, 0)

// NewDevice creates and initializes a new BLE112Device
// given a particular port
func NewDevice(port string) (*Device, error) {
	var device Device
	device.Port = port

	// stop scanning, clear buffer
	if err := device.Open(); err != nil {
		return nil, err
	}
	device.StopScan()
	device.Close()

	device.Open()
	device.MacAddress = device.GetAddress()
	device.Close()

	if device.MacAddress == nil || !bytes.Equal(device.MacAddress[3:], []byte{0x80, 0x07, 0x00}) {
		return nil, errors.New("Non-BLE112 MAC address detected")
	}
	return &device, nil
}

// Open opens the serial port connection to the BLE112
func (device *Device) Open() error {
	var err error
	c := serial.Config{Name: device.Port, Baud: 115200, ReadTimeout: time.Millisecond * 250}
	device.f, err = serial.OpenPort(&c)
	return err
}

// Close closes the serial port connection
func (device *Device) Close() {
	device.f.Close()
	device.f = nil
}

// SendCommand sends a command to a BLE112
func (device *Device) SendCommand(msgClass byte, msg byte, data []byte) (*Response, error) {
	dataSize := byte(len(data))
	cmd := []byte{BG_COMMAND, dataSize, msgClass, msg}
	cmd = append(cmd, data...)
	device.f.Write(cmd)

	return device.Read()
}

// GetAddress retrieves the BLE112's mac address.
func (device *Device) GetAddress() *beacon.MacAddress {
	var r *Response
	var err error
	retries := 4
	for err == nil && retries >= 0 {
		// sometimes it doesn't respond and we have to ask it again
		// not sure why.
		r, err = device.SendCommand(BG_MSG_CLASS_SYSTEM, BG_GET_ADDRESS, NULL_DATA)
		if err != nil {
			return nil
		}
		retries--
		if len(r.Data) >= 10 && bytes.Equal(r.Data[0:4], []byte{0, 6, 0, 2}) {
			break
		}
	}

	if err != nil || len(r.Data) < 10 {
		return nil
	}
	var macAddress beacon.MacAddress
	copy(macAddress[:], r.Data[4:10])
	return &macAddress
}

func (device *Device) StartAdvertising(data []byte) {
	device.SendCommand(BG_MSG_CLASS_CONNECTION, BG_DISCONNECT, NULL_DATA)
	device.SendCommand(BG_MSG_CLASS_GAP, BG_SET_MODE, []byte{BG_GAP_NON_DISCOVERABLE, BG_GAP_NON_CONNECTABLE})
	device.SendCommand(BG_MSG_CLASS_GAP, BG_GAP_SET_ADV_PARAM, []byte{0xA0, 0x00, 0xA0, 0x00, 0x07})
	b := append([]byte{0x00, byte(len(data) + 3), 0x02, 0x01, 0x06}, data...)
	device.SendCommand(BG_MSG_CLASS_GAP, BG_GAP_SET_ADV_DATA, b)
	device.SendCommand(BG_MSG_CLASS_GAP, BG_SET_MODE, []byte{BG_GAP_USER_DATA, BG_GAP_CONNECTABLE})
}

// advertiser.Advertiser interface

// AdvertiseMfgData advertises manufacturer data using the given mfg id
func (device *Device) AdvertiseMfgData(id uint16, ad advertiser.Advertisement) {
	device.Open()
	header := []byte{uint8(len(ad) + 1), 0xff, uint8(id), uint8(id >> 8)}
	device.StartAdvertising(append(header, ad[2:]...))
	device.Close()
}

// AdvertiseServiceData advertises the given service data with the given service uuid
func (device *Device) AdvertiseServiceData(id uint16, ad advertiser.Advertisement) {
	device.Open()
	header := []byte{
		0x03, 0x03, uint8(id), uint8(id >> 8),
		uint8(len(ad) + 1), 0x16, uint8(id), uint8(id >> 8),
	}
	device.StartAdvertising(append(header, ad[2:]...))
	device.Close()
}

// StopAdvertising stops advertising data
func (device *Device) StopAdvertising() {
	device.Open()
	device.SendCommand(BG_MSG_CLASS_GAP, BG_SET_MODE, []byte{BG_GAP_NON_DISCOVERABLE, BG_GAP_NON_CONNECTABLE})
	device.Close()
}

// StartScan tells the BLE112 to start scanning.
func (device *Device) StartScan() {
	device.SendCommand(BG_MSG_CLASS_CONNECTION, BG_DISCONNECT, NULL_DATA)
	device.SendCommand(BG_MSG_CLASS_GAP, BG_SET_MODE, []byte{BG_GAP_NON_DISCOVERABLE, BG_GAP_NON_CONNECTABLE})
	device.SendCommand(BG_MSG_CLASS_GAP, BG_DISCOVER_STOP, NULL_DATA)
	scanParams := []byte{200, 0, 200, 0, 0}
	device.SendCommand(BG_MSG_CLASS_GAP, BG_SCAN_PARAMS, scanParams)
	device.SendCommand(BG_MSG_CLASS_GAP, BG_DISCOVER, []byte{BG_GAP_DISCOVER_ALL})
}

// StopScan tells the BLE112 to stop scanning.
func (device *Device) StopScan() {
	device.f.Write([]byte{BG_COMMAND, 0, BG_MSG_CLASS_GAP, BG_DISCOVER_STOP})
}

// Scan uses the BLE112 device to scan for advertisements. It appends scans to
// the data channel, and exits when it recieves something on the done channel.
func (device *Device) Scan(data chan beacon.ScanData, done chan bool) {
	device.Open()
	device.StartScan()

	readChan := make(chan *Response, 2)
	shouldStop := false
	go func() {
		for !shouldStop {
			r, err := device.Read()
			if err == nil {
				readChan <- r
			}
		}
		close(readChan)
	}()

loop:
	for {
		select {
		case r, more := <-readChan:
			if !more {
				break loop
			}
			if r.IsAdvertisement() {
				if r.IsMfgAd() {
					data <- beacon.ScanData{r.Data[20:], r.MacAddress().String(), r.RSSI(), &r.Data}
				} else {
					data <- beacon.ScanData{r.Data[24:], r.MacAddress().String(), r.RSSI(), &r.Data}
				}
			}
		case <-done:
			shouldStop = true
		}
	}
	close(data)
	device.StopScan()
	device.Close()
}

// Read from the BLE112 device
func (device *Device) Read() (*Response, error) {
	var err error
	var byteCount int
	var output []byte

	if device.f == nil {
		return nil, errors.New("Device already closed!")
	}
	header := make([]byte, 4)
	byteCount, err = device.f.Read(header)
	if err != nil {
		return nil, err
	}
	output = append(output, header...)

	bytesLeft := int(header[1])
	for bytesLeft > 0 {
		buffer := make([]byte, bytesLeft)
		byteCount, err = device.f.Read(buffer)
		if err != nil {
			return nil, err
		}
		output = append(output, buffer[0:byteCount]...)
		bytesLeft -= byteCount
	}

	return &Response{output}, err
}

// DevicePaths returns a list of paths that correspond with possible
// BL112 devices.
func DevicePaths() ([]string, error) {
	if runtime.GOOS == "windows" {
		return []string{"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9"}, nil
	}

	paths, err := filepath.Glob("/dev/ttyACM*")
	if len(paths) == 0 {
		paths, err = filepath.Glob("/dev/cu.usbmodem*")
	}
	return paths, err
}

// Devices finds all the BLE112 devices that are currently on the system.
func Devices() []*Device {
	var devices []*Device
	paths, err := DevicePaths()
	check(err)
	for _, port := range paths {
		var device, _ = NewDevice(port)
		if device != nil && device.MacAddress != nil {
			devices = append(devices, device)
		}
	}
	return devices
}

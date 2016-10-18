package beacon

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// A BLE112Device represents a USB connected BLE112 which can be used for
// BLE scanning or advertising.
type BLE112Device struct {
	Port       string
	MacAddress *MacAddress
	f          *os.File
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
	BG_GAP_DISCOVER_ALL     = byte(2)
	BG_EVENT                = byte(0x80)
)

var NULL_DATA = make([]byte, 0)

var sttyCmdFormat = "-F %v 115200 raw -brkint -icrnl -imaxbel -opost -isig -icanon -iexten -echo -echoe -echok -echoctl -echoke"

// NewBLE112Device creates and initializes a new BLE112Device
// given a particular port
func NewBLE112Device(port string) *BLE112Device {

	if runtime.GOOS == "linux" {
		err := exec.Command("stty", fmt.Sprintf(sttyCmdFormat, port)).Run()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}

	var device BLE112Device
	device.Port = port

	// stop scanning, clear buffer
	device.Open()
	device.StopScan()
	device.Close()

	device.Open()
	device.MacAddress = device.GetAddress()
	device.Close()
	return &device
}

// Open opens the serial port connection to the BLE112
func (device *BLE112Device) Open() {
	device.f, _ = os.OpenFile(device.Port, os.O_RDWR, os.ModeDevice)
}

// Close closes the serial port connection
func (device *BLE112Device) Close() {
	device.f.Close()
	device.f = nil
}

// SendCommand sends a command to a BLE112
func (device *BLE112Device) SendCommand(msgClass byte, msg byte, data []byte) (*BLE112Response, error) {
	dataSize := byte(len(data))
	cmd := []byte{BG_COMMAND, dataSize, msgClass, msg}
	cmd = append(cmd, data...)
	device.f.Write(cmd)

	return device.Read()
}

// GetAddress retrieves the BLE112's mac address.
func (device *BLE112Device) GetAddress() *MacAddress {
	var r *BLE112Response
	var err error
	retries := 4
	for err == nil && retries >= 0 {
		// sometimes it doesn't respond and we have to ask it again
		// not sure why.
		r, err = device.SendCommand(BG_MSG_CLASS_SYSTEM, BG_GET_ADDRESS, NULL_DATA)
		retries--
		if len(r.Data) >= 10 && bytes.Equal(r.Data[0:4], []byte{0, 6, 0, 2}) {
			break
		}
	}

	if err != nil || len(r.Data) < 10 {
		return nil
	} else {
		var macAddress MacAddress
		copy(macAddress[:], r.Data[4:10])
		return &macAddress
	}
}

// StartScan tells the BLE112 to start scanning.
func (device *BLE112Device) StartScan() {
	device.SendCommand(BG_MSG_CLASS_CONNECTION, BG_DISCONNECT, NULL_DATA)
	device.SendCommand(BG_MSG_CLASS_GAP, BG_SET_MODE, []byte{BG_GAP_NON_DISCOVERABLE, BG_GAP_NON_CONNECTABLE})
	device.SendCommand(BG_MSG_CLASS_GAP, BG_DISCOVER_STOP, NULL_DATA)
	scanParams := []byte{200, 0, 200, 0, 0}
	device.SendCommand(BG_MSG_CLASS_GAP, BG_SCAN_PARAMS, scanParams)
	device.SendCommand(BG_MSG_CLASS_GAP, BG_DISCOVER, []byte{BG_GAP_DISCOVER_ALL})
}

// StopScan tells the BLE112 to stop scanning.
func (device *BLE112Device) StopScan() {
	device.f.Write([]byte{BG_COMMAND, 0, BG_MSG_CLASS_GAP, BG_DISCOVER_STOP})
}

// Scan uses the BLE112 device to scan for advertisements. It appends scans to
// the data channel, and exits when it recieves something on the done channel.
func (device *BLE112Device) Scan(data chan ScanData, done chan bool) {
	device.Open()
	device.StartScan()

	readChan := make(chan *BLE112Response, 2)
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
					data <- ScanData{r.Data[20:], r.MacAddress().String(), r.RSSI(), &r.Data}
				} else {
					data <- ScanData{r.Data[24:], r.MacAddress().String(), r.RSSI(), &r.Data}
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
func (device *BLE112Device) Read() (*BLE112Response, error) {
	var err error
	var byteCount int
	var output []byte

	if device.f == nil {
		return nil, errors.New("Device alerady closed!")
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

	var r BLE112Response
	r.Data = output
	return &r, err
}

// ble112DevicePaths returns a list of paths that correspond with possible
// BL112 devices.
func ble112DevicePaths() ([]string, error) {
	paths, err := filepath.Glob("/dev/ttyACM*")
	if len(paths) == 0 {
		paths, err = filepath.Glob("/dev/cu.usbmodem*")
	}
	return paths, err
}

// BLE112Devices finds all the BLE112 devices that are currently on the system.
func BLE112Devices() []*BLE112Device {
	var devices []*BLE112Device
	paths, err := ble112DevicePaths()
	check(err)
	for _, port := range paths {
		var device = NewBLE112Device(port)
		if device.MacAddress != nil {
			devices = append(devices, device)
		}
	}
	return devices
}

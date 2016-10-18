package beacon

import (
	"fmt"
	"sync"
	"time"
)

// a Scanner scans for beacons with a given ble interface (i.e., BlueZ, BLE112, CoreBluetooth)
type Scanner struct {
	device        ScanDevice
	parsers       []*Parser
	beaconChannel chan BeaconSlice
	done          chan bool
	beacons       BeaconSlice
}

// a ScanData represents a possible beacon advertisement that can be parsed into a beacon
type ScanData struct {
	Bytes  []byte
	Device string
	RSSI   int8
	Raw    *[]byte
}

// A ScanDevice will return ScanData on a channel.  Currently the only implementation is
// BLE112Device.
type ScanDevice interface {
	Scan(data chan ScanData, done chan bool)
}

// NewScanner initializes a new Scanner which will scan using the given ScanDevice and
// look for beacons given the list of beacon Parsers.
func NewScanner(d ScanDevice, p []*Parser) *Scanner {
	var s Scanner
	s.device = d
	s.parsers = p
	return &s
}

// Scan will scan for beacons and return a list of beacons that it detects on the interval
// given in cycleTime. It will stop scanning when it receives something on the done channel.
func (s *Scanner) Scan(cycleTime time.Duration, output chan BeaconSlice, done chan bool) {
	data := make(chan ScanData)
	doneOut := make(chan bool)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		timer := time.NewTimer(cycleTime)
		s.beacons = s.beacons[:0] // clear beacons slice
	loop:
		for {
			select {
			case scan, more := <-data:
				if !more {
					break loop
				}
				s.processScan(scan)
			case <-done:
				doneOut <- true
			case <-timer.C:
				output <- s.beacons
				s.beacons = s.beacons[:0] // clear beacons slice
				timer = time.NewTimer(cycleTime)
			}
		}
		fmt.Printf("boom\n")
		timer.Stop()
		wg.Done()
	}()
	s.device.Scan(data, doneOut)
	wg.Wait()
}

func (s *Scanner) processScan(scan ScanData) {
	beacon := Parse(scan.Bytes, s.parsers)
	if beacon == nil {
		return
	}
	beacon.Device = scan.Device
	found := s.beacons.Find(beacon)
	if found != nil {
		found.AddRSSI(scan.RSSI)
	} else {
		beacon.AddRSSI(scan.RSSI)
		s.beacons = append(s.beacons, beacon)
	}
}

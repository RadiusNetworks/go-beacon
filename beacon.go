// Package beacon provides functionality that allows interacting
// with bluetooth low-energy beacons.
package beacon

import (
	"fmt"
	"reflect"
	"strings"
)

// A Beacon represents a BLE beacon and is made up of IDs and various data.
type Beacon struct {
	Type   string
	Ids    Fields
	Data   Fields
	Power  Field
	rssis  []int8
	Device string
}

// A BeaconSlice is a list of Beacons
type BeaconSlice []*Beacon

// NewBeacon creates a new struct of type Beacon
func NewBeacon(t string, ids []Field, data []Field, power Field) Beacon {
	var beacon Beacon
	beacon.Type = t
	beacon.Ids = ids
	beacon.Data = data
	beacon.Power = power
	return beacon
}

// Generates a description of the beacon
func (b *Beacon) String() string {
	idStrings := make([]string, len(b.Ids))
	for i, id := range b.Ids {
		idStrings[i] = id.String()
	}
	idString := strings.Join(idStrings, " ")

	return fmt.Sprintf("%v - %v: %v, rssi: %.2f, scans: %v", b.Device, b.Type, idString, b.RSSI(), len(b.rssis))
}

func (b *Beacon) RSSI() float64 {
	total := 0.0
	for _, rssi := range b.rssis {
		total += float64(rssi)
	}
	return total / float64(len(b.rssis))
}

func (b *Beacon) AddRSSI(rssi int8) {
	b.rssis = append(b.rssis, rssi)
}

func (a *Beacon) Equal(b *Beacon) bool {
	return strings.Compare(a.Device, b.Device) == 0 && reflect.DeepEqual(a.Ids, b.Ids)
}

func (beacons BeaconSlice) Find(b *Beacon) *Beacon {
	for _, item := range beacons {
		if b.Equal(item) {
			return item
		}
	}
	return nil
}

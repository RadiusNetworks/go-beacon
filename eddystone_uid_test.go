package beacon

import "testing"

func TestNewEddystoneUIDBeacon(t *testing.T) {
	namespace := "00010203040506070809"
	instance := "00010203040506"
	beacon, error := NewEddystoneUIDBeacon(namespace, instance, -42)
	if error != nil {
		t.Error("Unable to create Eddystone UID beacon")
	}
	if beacon.Type != BeaconTypeEddystoneUID {
		t.Errorf("Beacon type %v; expected %v", beacon.Type, BeaconTypeEddystoneUID)
	}
}

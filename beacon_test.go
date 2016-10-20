package beacon

import (
	"strings"
	"testing"
)

var (
	uuidString      = "66484d6e54bf4d67b2698b100151510b"
	anyUuid         = FieldFromHex(uuidString)
	anyMajor        = FieldFromUint16(1)
	anyMinor        = FieldFromUint16(5)
	anyAltbeaconIds = []Field{anyUuid, anyMajor, anyMinor}
)

func TestNewBeacon(t *testing.T) {
	beacon := NewBeacon(
		"altbeacon",
		anyAltbeaconIds,
		[]Field{},
		FieldFromInt8(-65),
	)

	if !beacon.Ids.Equal(anyAltbeaconIds) {
		t.Errorf("got %v; expected %v", beacon.Ids, anyAltbeaconIds)
	}
}

func TestString(t *testing.T) {
	beacon := NewBeacon(
		"altbeacon",
		anyAltbeaconIds,
		[]Field{},
		FieldFromInt8(-65),
	)

	s := beacon.String()
	if !strings.Contains(s, uuidString) {
		t.Errorf("expected \"%v\" to include the uuid.", s)
	}
}

func TestMacAddress(t *testing.T) {
	expected := "00:07:80:14:47:d5"
	ma := MacAddress{0xd5, 0x47, 0x14, 0x80, 0x07, 0x00}
	got := ma.String()
	if got != expected {
		t.Errorf("got %v; expected %v", got, expected)
	}
}

package beacon

import (
	"testing"
)

// malformedAltbeacon is an altbeacon advertisement with the last byte truncated
var malformedAltbeacon = "BEACBEACE858FC8A372B4BEFA05393F98CD4E1770001000140"

// properAltbeacon is a properly formed altbeacon advertisement
var properAltbeacon = "BEACBEACE858FC8A372B4BEFA05393F98CD4E177000100014020"

func TestParserMalformedAltbeacon(t *testing.T) {
	beaconBytes := []byte(FieldFromHex(malformedAltbeacon))
	beaconBytes = collapseCapacity(beaconBytes)
	beacon := Parse(beaconBytes, DefaultParsers())
	if beacon != nil {
		t.Errorf("Expected malformed beacon to not parse, but got a beacon: %v", beacon)
	}
}

func TestParserProperAltbeacon(t *testing.T) {
	beaconBytes := []byte(FieldFromHex(properAltbeacon))
	beaconBytes = collapseCapacity(beaconBytes)
	beacon := Parse(beaconBytes, DefaultParsers())
	if beacon == nil {
		t.Error("Expected beacon bytes to parse, but got nil")
	}
}

func collapseCapacity(in []byte) []byte {
	out := make([]byte, len(in))
	copy(out, in)
	return out
}

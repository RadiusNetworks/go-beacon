package beacon

// NewAltBeacon returns an altbeacon beacon.
func NewAltBeacon(uuid string, major uint16, minor uint16, pwr int8) *Beacon {
	beaconIds := AltBeaconFields(uuid, major, minor)
	beacon := NewBeacon("altbeacon",
		beaconIds,                   // ids
		Fields{FieldFromInt8(0x20)}, // data
		FieldFromInt8(pwr),          // measured power
	)
	return &beacon
}

// AltBeaconFields returns a slice containing the uuid, major, and minor.
func AltBeaconFields(uuid string, major uint16, minor uint16) Fields {
	return Fields{FieldFromHex(uuid), FieldFromUint16(major), FieldFromUint16(minor)}
}

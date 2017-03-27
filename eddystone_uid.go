package beacon

// BeaconTypeEddystoneUID indicates a beacon of type Eddystone-UID.
const BeaconTypeEddystoneUID = "eddystone_uid"

// NewEddystoneUIDBeacon returns an Eddystone-UID beacon or an error if
// the namespace or instance are invalid hex strings or the wrong length.
func NewEddystoneUIDBeacon(namespace string, instance string, pwr int8) (*Beacon, error) {
	beaconIds := EddystoneUIDFields(namespace, instance)

	beacon := NewBeacon(BeaconTypeEddystoneUID,
		beaconIds,          // ids
		Fields{},           // data
		FieldFromInt8(pwr), // measured power
	)
	return &beacon, nil
}

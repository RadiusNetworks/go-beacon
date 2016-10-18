package beacon

import "encoding/hex"

// A BLE112Response is data that the BLE112 returns while scanning or
// in response to a command.
type BLE112Response struct {
	Data []byte
}

func (r *BLE112Response) IsEvent() bool {
	return r.Data[0] == BG_EVENT
}

func (r *BLE112Response) IsGapScan() bool {
	return r.Data[2] == BG_MSG_CLASS_GAP && r.Data[3] == byte(0)
}

func (r *BLE112Response) IsMfgAd() bool {
	return len(r.Data) > 20 && r.Data[19] == byte(0xff)
}

func (r *BLE112Response) IsServiceAd() bool {
	return len(r.Data) > 24 && r.Data[19] == byte(0x03)
}

func (r *BLE112Response) IsAdvertisement() bool {
	return len(r.Data) > 20 && r.IsEvent() && r.IsGapScan() && (r.IsMfgAd() || r.IsServiceAd())
}

func (r *BLE112Response) MacAddress() *MacAddress {
	var a MacAddress
	copy(a[:], r.Data[6:12])
	return &a
}

func (r *BLE112Response) RSSI() int8 {
	return int8(r.Data[4])
}

func (r *BLE112Response) String() string {
	return hex.EncodeToString(r.Data)
}

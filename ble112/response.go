package ble112

import (
	"encoding/hex"
	"github.com/RadiusNetworks/go-beacon"
)

// A Response is data that the BLE112 returns while scanning or
// in response to a command.
type Response struct {
	Data []byte
}

func (r *Response) IsEvent() bool {
	return r.Data[0] == BG_EVENT
}

func (r *Response) IsGapScan() bool {
	return r.Data[2] == BG_MSG_CLASS_GAP && r.Data[3] == byte(0)
}

func (r *Response) IsMfgAd() bool {
	return len(r.Data) > 20 && r.Data[19] == byte(0xff)
}

func (r *Response) IsServiceAd() bool {
	return len(r.Data) > 24 && r.Data[19] == byte(0x03)
}

func (r *Response) IsAdvertisement() bool {
	return len(r.Data) > 20 && r.IsEvent() && r.IsGapScan() && (r.IsMfgAd() || r.IsServiceAd())
}

func (r *Response) AdData() []byte {
	if r.IsMfgAd() {
		return r.Data[20:]
	} else if r.IsServiceAd() {
		return r.Data[24:]
	} else {
		return []byte{}
	}
}

func (r *Response) MacAddress() *beacon.MacAddress {
	var a beacon.MacAddress
	copy(a[:], r.Data[6:12])
	return &a
}

func (r *Response) RSSI() int8 {
	return int8(r.Data[4])
}

func (r *Response) String() string {
	return hex.EncodeToString(r.Data)
}

package beacon

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// A MacAddress stores a BLE device's mac address.
type MacAddress [6]byte

// String converts a MacAddress into a hex string with bytes separated
// by colons.
func (addr MacAddress) String() string {
	return fmt.Sprintf(
		"%02x:%02x:%02x:%02x:%02x:%02x",
		addr[5], addr[4], addr[3], addr[2], addr[1], addr[0],
	)
}

// ParseMacAddress parses a MacAddress struct from a string with the
// format "00:11:22:33:44:55"
func ParseMacAddress(s string) MacAddress {
	bytes, _ := hex.DecodeString(strings.Replace(s, ":", "", -1))
	return MacAddress{bytes[5], bytes[4], bytes[3], bytes[2], bytes[1], bytes[0]}
}

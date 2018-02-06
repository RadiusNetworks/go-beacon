package beacon

import (
	"encoding/hex"
	"errors"
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

// MarshalJSON implements the json.Marshaler interface by returning the
// mac address as a string.
func (addr MacAddress) MarshalJSON() (text []byte, err error) {
	str := fmt.Sprintf("\"%v\"", addr.String())
	return []byte(str), nil
}

// UnmarshalJSON implements the json.Unmarshaler inferface by parsing
// the string representing the mac address.
func (addr *MacAddress) UnmarshalJSON(text []byte) error {
	str := string(text)
	if len(str) == 19 && str[0] == '"' {
		*addr = ParseMacAddress(string(text[1 : len(text)-1]))
	} else {
		return errors.New("\"%v\" is not a mac address")
	}
	return nil
}

// ParseMacAddress parses a MacAddress struct from a string with the
// format "00:11:22:33:44:55"
func ParseMacAddress(s string) MacAddress {
	bytes, _ := hex.DecodeString(strings.Replace(s, ":", "", -1))
	return MacAddress{bytes[5], bytes[4], bytes[3], bytes[2], bytes[1], bytes[0]}
}

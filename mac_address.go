package beacon

import "fmt"

// A MacAddress stores a BLE device's mac address.
type MacAddress [6]byte

// String converts a MacAddress into a hex string with bytes separated
// by colons.
func (addr *MacAddress) String() string {
	return fmt.Sprintf(
		"%02x:%02x:%02x:%02x:%02x:%02x",
		addr[5], addr[4], addr[3], addr[2], addr[1], addr[0],
	)
}

package beacon

import "encoding/hex"
import "strconv"
import "encoding/binary"

// A Field stores a beacon ID or data field.
type Field []byte

// String converts a beacon Field into a human readable format.
func (f *Field) String() string {
	switch len(*f) {
	case 1:
		return strconv.Itoa(int((*f)[0]))
	case 2:
		return strconv.Itoa(int(binary.BigEndian.Uint16((*f)[:])))
	default:
		return hex.EncodeToString(*f)
	}
}

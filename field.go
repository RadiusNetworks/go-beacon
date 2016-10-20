package beacon

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"strconv"
)

// A Field stores a beacon ID or data field.
type Field []byte
type Fields []Field

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

func FieldFromUint16(n uint16) Field {
	var field Field = make([]byte, 2)
	binary.BigEndian.PutUint16(field, n)
	return field
}

func FieldFromInt8(n int8) Field {
	var field Field = []byte{byte(n)}
	return field
}

func FieldFromHex(s string) Field {
	var field Field
	field, _ = hex.DecodeString(s)
	return field
}

func (a Field) Equal(b Field) bool {
	return bytes.Equal(a, b)
}

func (a Fields) Equal(b Fields) bool {
	for i, _ := range a {
		if !bytes.Equal(a[i], b[i]) {
			return false
		}
	}
	return true
}

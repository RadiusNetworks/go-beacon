package beacon

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// A Field stores a beacon ID or data field.
type Field []byte
type Fields []Field

// String converts a beacon Field into a human readable format.
func (f *Field) String() string {
	switch len(*f) {
	case 1:
		return strconv.Itoa(int(f.Int8()))
	case 2:
		return strconv.Itoa(int(f.Uint16()))
	default:
		return f.Hex()
	}
}

func (f Field) MarshalJSON() (text []byte, err error) {
	if len(f) == 2 {
		text = []byte(f.String())
	} else {
		text = []byte(fmt.Sprintf("\"%v\"", f.Hex()))
	}
	err = nil
	return
}

func (f *Field) Int8() int8 {
	return int8((*f)[0])
}

func (f *Field) Uint8() uint8 {
	return uint8((*f)[0])
}

func (f *Field) Uint16() uint16 {
	return binary.BigEndian.Uint16((*f)[:])
}

func (f *Field) Hex() string {
	return hex.EncodeToString(*f)
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
	// remove dashes in case we were given a UUID
	s = strings.Replace(s, "-", "", -1)
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

func (f *Fields) Hex() string {
	strFields := make([]string, len(*f))
	for i, field := range *f {
		strFields[i] = field.Hex()
	}
	return strings.Join(strFields, " ")
}

func (f *Fields) String() string {
	strFields := make([]string, len(*f))
	for i, field := range *f {
		strFields[i] = field.String()
	}
	return strings.Join(strFields, " ")
}

func (f *Fields) Key() string {
	data := make([][]byte, len(*f))
	for i, field := range *f {
		data[i] = field
	}
	return string(bytes.Join(data, []byte{}))
}

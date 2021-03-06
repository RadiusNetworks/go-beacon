package beacon

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// A Field stores a beacon ID or data field.
type Field []byte

// Fields is a slice of Field structs
type Fields []Field

// String converts a beacon Field into a human readable format.
func (f Field) String() string {
	switch len(f) {
	case 1:
		return strconv.Itoa(int(f.Int8()))
	case 2:
		return strconv.Itoa(int(f.Uint16()))
	default:
		return f.Hex()
	}
}

// MarshalJSON outputs a hex string for most fields, or an integer
// for fields that are two bytes in length.
func (f Field) MarshalJSON() (text []byte, err error) {
	if len(f) == 2 {
		text = []byte(f.String())
	} else {
		text = []byte(fmt.Sprintf("\"%v\"", f.Hex()))
	}
	err = nil
	return
}

// UnmarshalJSON unmarshals a json string into a Field. If the text starts with
// a quote, it assumes it is a hex string. Otherwise it will unmarshal it as a
// uint16 field, if it looks like an integer.
func (f *Field) UnmarshalJSON(text []byte) error {
	str := string(text)
	if str[0] == '"' {
		*f = FieldFromHex(str[1 : len(str)-1])
	} else if i, err := strconv.Atoi(str); err == nil && strconv.Itoa(i) == str {
		*f = FieldFromUint16(uint16(i))
	} else {
		return errors.New("Error unmarshaling text - invalid value")
	}
	return nil
}

// Int8 returns the field as an int8 value
func (f *Field) Int8() int8 {
	return int8((*f)[0])
}

// Uint8 returns the field as a uint8 value
func (f *Field) Uint8() uint8 {
	return uint8((*f)[0])
}

// Uint16 returns the field as a uint16 value. This is especially helpful for
// ibeacon or altbeaon major and minor fields.
func (f *Field) Uint16() uint16 {
	return binary.BigEndian.Uint16((*f)[:])
}

// Hex returns the field as a hex string
func (f *Field) Hex() string {
	return hex.EncodeToString(*f)
}

// FieldFromUint16 converts a uint16 value into a Field
func FieldFromUint16(n uint16) Field {
	var field Field = make([]byte, 2)
	binary.BigEndian.PutUint16(field, n)
	return field
}

// FieldFromInt8 converts an int8 value into a Field struct
func FieldFromInt8(n int8) Field {
	var field Field = []byte{byte(n)}
	return field
}

// FieldFromHex converts a hex string into a Field struct
func FieldFromHex(s string) Field {
	var field Field
	// remove dashes in case we were given a UUID
	s = strings.Replace(s, "-", "", -1)
	field, _ = hex.DecodeString(s)
	return field
}

// Equal tests if two field structs have the same bytes
func (a Field) Equal(b Field) bool {
	return bytes.Equal(a, b)
}

// Equal tests if two Field slices contain matching Field structs
func (a Fields) Equal(b Fields) bool {
	for i := range a {
		if !bytes.Equal(a[i], b[i]) {
			return false
		}
	}
	return true
}

// UUIDMajorMinorFields returns Fields suitable for an altBeacon or iBeacon
func UUIDMajorMinorFields(uuid string, major uint16, minor uint16) Fields {
	return Fields{
		FieldFromHex(uuid),
		FieldFromUint16(major),
		FieldFromUint16(minor),
	}
}

// EddystoneUIDFields returns Fields suitable for an Eddystone UID beacon
func EddystoneUIDFields(namespace, instance string) Fields {
	return Fields{
		FieldFromHex(namespace),
		FieldFromHex(instance),
	}
}

// Hex returns hex strings for each Field separated by a space
func (f *Fields) Hex() string {
	strFields := make([]string, len(*f))
	for i, field := range *f {
		strFields[i] = field.Hex()
	}
	return strings.Join(strFields, " ")
}

// String returns a human readable representation of Fields
func (f *Fields) String() string {
	strFields := make([]string, len(*f))
	for i, field := range *f {
		strFields[i] = field.String()
	}
	return strings.Join(strFields, " ")
}

// Key returns a value which can be used as a map key to uniquely
// represent this set of fields
func (f *Fields) Key() string {
	data := make([][]byte, len(*f))
	for i, field := range *f {
		data[i] = field
	}
	return string(bytes.Join(data, []byte{}))
}

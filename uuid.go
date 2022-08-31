package beacon

import (
	"bytes"
	"encoding/binary"
)

// UUID represents a 16 byte UUID
type UUID [16]byte

func (uuid UUID) ReversedUUID() []byte {
	var a UUID
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = uuid[opp], uuid[i]
	}
	return a[:]
}

func (uuid *UUID) Reverse() {
	for i, j := 0, len(uuid)-1; i < j; i, j = i+1, j-1 {
		uuid[i], uuid[j] = uuid[j], uuid[i]
	}
}

// MarshalJSON implements json.Marshaler
func (uuid UUID) MarshalJSON() (text []byte, err error) {
	return Field(uuid[:]).MarshalJSON()
}

// UnmarshalJSON implmenets json.Unmarshaler
func (uuid *UUID) UnmarshalJSON(text []byte) error {
	var field Field
	if err := field.UnmarshalJSON(text); err != nil {
		return err
	}
	copy((*uuid)[:], field)
	return nil
}

func (uuid UUID) AddIndex(i uint32) UUID {
	var newUUID UUID
	var firstFour uint32
	reader := bytes.NewReader(uuid[0:4])
	binary.Read(reader, binary.BigEndian, &firstFour)

	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, uint32(firstFour+i))
	copy(newUUID[0:4], buf.Bytes())
	copy(newUUID[4:], uuid[4:])
	return newUUID
}

package beacon

import (
	"reflect"
	"testing"
)

func TestUUIDReverse(t *testing.T) {
	uuid := UUID{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p'}
	reverse := []byte{'p', 'o', 'n', 'm', 'l', 'k', 'j', 'i', 'h', 'g', 'f', 'e', 'd', 'c', 'b', 'a'}
	if !reflect.DeepEqual(uuid.ReversedUUID(), reverse) {
		t.Errorf("got %v; expected %v", uuid.ReversedUUID(), reverse)
	}

	uuid.Reverse()
	reverseUUID := UUID{'p', 'o', 'n', 'm', 'l', 'k', 'j', 'i', 'h', 'g', 'f', 'e', 'd', 'c', 'b', 'a'}
	if !reflect.DeepEqual(uuid, reverseUUID) {
		t.Errorf("got %v; expected %v", uuid, reverseUUID)
	}
}

func testEq(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

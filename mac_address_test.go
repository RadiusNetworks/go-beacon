package beacon

import (
	"testing"
)

type mac struct {
	hex  string
	bits [6]byte
	text []byte
}

var testCases []mac

func init() {
	testCases = []mac{
		{
			"fe:fd:55:c9:64:24",
			[6]byte{0x24, 0x64, 0xc9, 0x55, 0xfd, 0xfe},
			[]byte{'"', 'f', 'e', ':', 'f', 'd', ':', '5', '5', ':', 'c', '9', ':', '6', '4', ':', '2', '4', '"'},
		},
		{
			"29:fe:fd:41:21:20",
			[6]byte{0x20, 0x21, 0x41, 0xfd, 0xfe, 0x29},
			[]byte{'"', '2', '9', ':', 'f', 'e', ':', 'f', 'd', ':', '4', '1', ':', '2', '1', ':', '2', '0', '"'},
		},
	}
}

func TestMacString(t *testing.T) {
	for _, tc := range testCases {
		result := MacAddress(tc.bits).String()
		if result != tc.hex {
			t.Errorf("result %s; want %s", result, tc.hex)
		}
	}
}

func TestMacMarshalJSON(t *testing.T) {
	for _, tc := range testCases {
		result, _ := MacAddress(tc.bits).MarshalJSON()
		if !lenEqual(result, tc.text) {
			t.Errorf("result %s; want %s", result, tc.text)
		}
		if !valueEqual(result, tc.text) {
			t.Errorf("result %s; want %s", result, tc.text)
		}
	}
}

func TestMacUnmarshalJSON(t *testing.T) {
	for _, tc := range testCases {
		addr := MacAddress{}
		addr.UnmarshalJSON(tc.text)
		if addr != tc.bits {
			t.Errorf("result %s; want %v", addr, tc.hex)
		}
	}
}

func TestParseMacAddress(t *testing.T) {
	for _, tc := range testCases {
		want := MacAddress(tc.bits)
		result := ParseMacAddress(tc.hex)
		if result != want {
			t.Errorf("result %s; want %s", result, want)
		}
	}
}

func lenEqual(r []byte, t []byte) bool {
	return len(r) == len(t)
}

func valueEqual(r []byte, t []byte) (equal bool) {
	for i, v := range r {
		if v == t[i] {
			equal = true
		} else {
			equal = false
			break
		}
	}
	return equal
}

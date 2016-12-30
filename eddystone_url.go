package beacon

import (
	"bytes"
	"fmt"
)

type scheme struct {
	name string
	code byte
}

var (
	schemes = []scheme{
		scheme{"http://www.", 0x00},
		scheme{"https://www.", 0x01},
		scheme{"http://", 0x02},
		scheme{"https://", 0x03},
	}

	expansions = map[string]byte{
		".com/":  0x00,
		".org/":  0x01,
		".edu/":  0x02,
		".net/":  0x03,
		".info/": 0x04,
		".biz/":  0x05,
		".gov/":  0x06,
		".com":   0x07,
		".org":   0x08,
		".edu":   0x09,
		".net":   0x0a,
		".info":  0x0b,
		".biz":   0x0c,
		".gov":   0x0d,
	}
)

// NewEddystoneURLBeacon returns an Eddystone-URL beacon or an error if
// the URL cannot be compressed.
func NewEddystoneURLBeacon(url string, pwr int8) (*Beacon, error) {
	beaconIds, err := EddystoneURLFields(url)
	if err != nil {
		return nil, err
	}

	beacon := NewBeacon("eddystone_url",
		beaconIds,          // ids
		Fields{},           // data
		FieldFromInt8(pwr), // measured power
	)
	return &beacon, nil
}

// EddystoneURLFields returns a slice (length 1) with the compressed url
func EddystoneURLFields(url string) (Fields, error) {
	f, err := CompressEddystoneURL(url)
	return Fields{f}, err
}

// CompressEddystoneURL compresses a url as defined in the Eddystone URL spec
func CompressEddystoneURL(url string) (Field, error) {
	urlBytes := []byte(url)
	foundScheme := false
	for _, scheme := range schemes {
		schemeBytes := []byte(scheme.name)
		if bytes.HasPrefix(urlBytes, schemeBytes) {
			urlBytes = bytes.Replace(urlBytes, schemeBytes, []byte{scheme.code}, 1)
			foundScheme = true
			break
		}
	}

	if !foundScheme {
		return urlBytes, fmt.Errorf("URL does not have valid scheme")
	}

	for expansion, code := range expansions {
		expansionBytes := []byte(expansion)
		if bytes.Contains(urlBytes, expansionBytes) {
			urlBytes = bytes.Replace(urlBytes, expansionBytes, []byte{code}, -1)
		}
	}

	if len(urlBytes) > 18 {
		return urlBytes, fmt.Errorf("URL is too long")
	} else {
		return urlBytes, nil
	}
}

// DecompressEddystoneURL decompresses a Field into a url, as defined in the
// Eddystene URL spec
func (f *Field) DecompressEddystoneURL() string {
	url := ""
	for _, scheme := range schemes {
		if (*f)[0] == scheme.code {
			url += scheme.name
			break
		}
	}

	urlBytes := (*f)[1:]
	for expansion, code := range expansions {
		if bytes.Contains(urlBytes, []byte{code}) {
			urlBytes = bytes.Replace(urlBytes, []byte{code}, []byte(expansion), -1)
		}
	}
	url += string(urlBytes)

	return url
}

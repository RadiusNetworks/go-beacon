package beacon

import "strings"
import "encoding/hex"
import "strconv"
import "bytes"

var DefaultLayouts = map[string]string{
	"altbeacon":     "m:2-3=beac,i:4-19,i:20-21,i:22-23,p:24-24,d:25-25",
	"eddystone_uid": "s:0-1=feaa,m:2-2=00,p:3-3:-41,i:4-13,i:14-19,d:20-21",
	"eddystone_url": "s:0-1=feaa,m:2-2=10,p:3-3:-41,i:4-21v",
	"eddystone_tlm": "s:0-1=feaa,m:2-2=20,d:3-3,d:4-5,d:6-7,d:8-11,d:12-15",
	"eddystone_eid": "s:0-1=feaa,m:2-2=30,p:3-3:-41,i:4-11",
}

// DefaultParsers returns a list of beacon parsers defined by default.
func DefaultParsers() []*Parser {
	DefaultParser := make([]*Parser, len(DefaultLayouts))
	i := 0
	for name, layout := range DefaultLayouts {
		DefaultParser[i] = NewParser(name, layout)
		i++
	}
	return DefaultParser
}

type fieldParams struct {
	start     int
	end       int
	length    int
	varLength bool
	expected  []byte
}

// A Parser can parse beacon advertisements.
type Parser struct {
	Name       string
	Layout     string
	matchers   []fieldParams
	idFields   []fieldParams
	dataFields []fieldParams
	powerField fieldParams
}

// NewParser initializes a new beacon parser with the given name and layout.
func NewParser(name string, layout string) *Parser {
	var p Parser
	p.Name = name
	p.Layout = layout
	p.parseLayout(layout)
	return &p
}

func (p *Parser) parseLayout(layout string) {
	parts := strings.Split(layout, ",")
	for _, part := range parts {
		var params fieldParams
		details := strings.Split(part, ":")
		partType := details[0]
		details = strings.Split(details[1], "=")
		startEnd := strings.Split(details[0], "-")
		params.varLength = strings.HasSuffix(details[0], "v")
		params.start, _ = strconv.Atoi(startEnd[0])
		params.end, _ = strconv.Atoi(startEnd[1])
		params.length = params.end - params.start + 1
		if len(details) > 1 {
			params.expected, _ = hex.DecodeString(details[1])
		}
		switch partType {
		case "m":
			p.matchers = append(p.matchers, params)
		case "s":
			// swap bytes for service UUID
			params.expected = []byte{params.expected[1], params.expected[0]}
			p.matchers = append(p.matchers, params)
		case "i":
			p.idFields = append(p.idFields, params)
		case "d":
			p.dataFields = append(p.dataFields, params)
		case "p":
			p.powerField = params
		}
	}
}

// Matches returns true if the advertisement data matches this layout.
func (p *Parser) Matches(data []byte) bool {
	for _, params := range p.matchers {
		if !bytes.Equal(data[params.start:params.end+1], params.expected) {
			return false
		}
	}
	return true
}

// ParseIds parses a beacon's IDs out of advertisement data, according
// to the layout.
func (p *Parser) ParseIds(data []byte) []Field {
	var ids = make([]Field, len(p.idFields))
	for i, params := range p.idFields {
		var id []byte
		if params.varLength {
			id = data[params.start:]
		} else {
			id = data[params.start : params.end+1]
		}
		ids[i] = id
	}
	return ids
}

// ParseData parses a beacon's data fields out of advertisement data, according
// to the layout.
func (p *Parser) ParseData(data []byte) []Field {
	var fields = make([]Field, len(p.dataFields))
	for i, params := range p.dataFields {
		var field []byte
		if params.varLength {
			field = data[params.start:]
		} else {
			field = data[params.start : params.end+1]
		}
		fields[i] = field
	}
	return fields
}

// ParsePower parses a beacon's measured power field out of advertisement data, according
// to the layout.
func (p *Parser) ParsePower(data []byte) Field {
	return data[p.powerField.start : p.powerField.end+1]
}

// Parse parses advertisement data. It returns an instance of Beacon if it
// matches the layout, otherwise it returns nil.
func (p *Parser) Parse(data []byte) *Beacon {
	if !p.Matches(data) {
		return nil
	}
	ids := p.ParseIds(data)
	dataFields := p.ParseData(data)
	measuredPower := p.ParsePower(data)
	beacon := NewBeacon(p.Name, ids, dataFields, measuredPower)
	return &beacon
}

// GenerateAd generates the bytes of a beacon advertisement with the
// given beacon.
func (p *Parser) GenerateAd(b *Beacon) []byte {
	insertField := func(ad *[]byte, p fieldParams, field Field) {
		l := len(field)
		end := p.start + l
		if len(*ad) < end {
			*ad = (*ad)[:end]
		}
		copy((*ad)[p.start:p.start+l], field)
	}

	ad := make([]byte, 0, 32)
	for idx, field := range p.idFields {
		insertField(&ad, field, b.Ids[idx])
	}
	for idx, field := range p.dataFields {
		insertField(&ad, field, b.Data[idx])
	}
	insertField(&ad, p.powerField, b.Power)
	for _, field := range p.matchers {
		insertField(&ad, field, field.expected)
	}
	return ad
}

// Parse attempts to parse a Beacon from advertisement data, given a list
// of Parsers.
func Parse(data []byte, parsers []*Parser) *Beacon {
	for _, parser := range parsers {
		beacon := parser.Parse(data)
		if beacon != nil {
			return beacon
		}
	}
	return nil
}

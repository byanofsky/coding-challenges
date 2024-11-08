package internal

import (
	"fmt"
	"regexp"
)

var BULK_STRING_LEN = regexp.MustCompile(`^\$(-?\d+)\r\n`)

func Deserialize(s string) (*Data, error) {
	f := s[0]
	switch f {
	case '$':
		return deserializeBulkStringOrNull(s)
	default:
		return nil, fmt.Errorf("error unexpected first char: %v", f)
	}
}

func deserializeBulkStringOrNull(s string) (*Data, error) {
	m := BULK_STRING_LEN.FindStringSubmatch(s)

	// Length
	l := m[1]

	// Null
	if l == "-1" {
		return &Data{kind: NullKind}, nil
	}

	fmt.Println(m)
	return nil, nil
}

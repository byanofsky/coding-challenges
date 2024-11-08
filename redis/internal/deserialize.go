package internal

import (
	"fmt"
	"regexp"
)

var BULK_STRING_LEN = regexp.MustCompile(`^\$(-?\d+)\r\n`)
var SIMPLE_STRING_LEN = regexp.MustCompile(`^\+([^\n\r]+)\r\n`)

func Deserialize(s string) (*Data, error) {
	f := s[0]
	switch f {
	case '$':
		return deserializeBulkStringOrNull(s)
	case '+':
		return deserializeSimpleString(s)
	default:
		return nil, fmt.Errorf("error unexpected first char: %v", f)
	}
}

func deserializeBulkStringOrNull(s string) (*Data, error) {
	m := BULK_STRING_LEN.FindStringSubmatch(s)
	if m == nil {
		return nil, fmt.Errorf("error bulk string format: %q", s)
	}

	// Length
	l := m[1]

	// Null
	if l == "-1" {
		return &Data{kind: NullKind}, nil
	}

	fmt.Println(m)
	return nil, nil
}

func deserializeSimpleString(s string) (*Data, error) {
	m := SIMPLE_STRING_LEN.FindStringSubmatch(s)
	if m == nil {
		return nil, fmt.Errorf("error simple string format: %q", s)
	}

	return &Data{kind: SimpleStringKind, value: m[1]}, nil
}

package internal

import (
	"fmt"
	"regexp"
	"strconv"
)

var BULK_STRING_RE = regexp.MustCompile(`(?s)^(-?\d+)\r\n(.*)$`)
var SIMPLE_STRING_RE = regexp.MustCompile(`^([^\n\r]+)\r\n(.*)$`)
var INT_RE = regexp.MustCompile(`^((?:\+|-)?\d+)\r\n(.*)$`)
var ARRAY_RE = regexp.MustCompile(`(?s)^(\d+)\r\n(.*)$`)
var SIMPLE_ERROR_RE = regexp.MustCompile(`(?s)^([^\n\r]+)\r\n(.*)$`)

func Deserialize(s string) (*Data, error) {
	d, remaining, err := parseDeserialize(s)
	if len(remaining) != 0 {
		return nil, fmt.Errorf("error RESP format remaining: %q", remaining)
	}
	return d, err
}

func deserializeBulkStringOrNull(s string) (*Data, string, error) {
	m := BULK_STRING_RE.FindStringSubmatch(s)
	if m == nil {
		return nil, "", fmt.Errorf("error bulk string format: %q", s)
	}

	// Length
	l, err := strconv.Atoi(m[1])
	if err != nil {
		return nil, "", fmt.Errorf("error parsing length: %q", s)
	}

	if l < -1 {
		return nil, "", fmt.Errorf("error invalid length: %d %q", l, s)
	}

	// Null
	if l == -1 {
		return &Data{kind: NullKind}, "", nil
	}

	// Bulk String
	match := m[2]
	// Bulk string value
	value := match[0:l]
	// Should bulk string must end with CRLF
	if match[l:l+2] != "\r\n" {
		return nil, "", fmt.Errorf("error bulk string format, should end with CRLF: %q", match[0:l+2])
	}
	remaining := match[l+2:]
	return &Data{kind: BulkStringKind, value: value}, remaining, nil
}

func deserializeSimpleString(s string) (*Data, string, error) {
	m := SIMPLE_STRING_RE.FindStringSubmatch(s)
	if m == nil {
		return nil, "", fmt.Errorf("error simple string format: %q", s)
	}
	match := m[1]
	remaining := m[2]

	return &Data{kind: SimpleStringKind, value: match}, remaining, nil
}

func deserializeInt(s string) (*Data, string, error) {
	m := INT_RE.FindStringSubmatch(s)
	if m == nil {
		return nil, "", fmt.Errorf("error int format: %q", s)
	}
	match := m[1]
	remaining := m[2]

	i, err := strconv.Atoi(match)
	if err != nil {
		return nil, "", fmt.Errorf("error converting int %s: %w", m[1], err)
	}

	return &Data{kind: IntKind, value: i}, remaining, nil
}

func deserializeArray(s string) (*Data, string, error) {
	var remaining string
	m := ARRAY_RE.FindStringSubmatch(s)
	if m == nil {
		return nil, "", fmt.Errorf("error array format: %q", s)
	}

	length, err := strconv.Atoi(m[1])
	if err != nil {
		return nil, remaining, fmt.Errorf("error parsing array length: %w", err)
	}
	remaining = m[2]
	value := make([]*Data, 0, length)

	for length > 0 {
		var element *Data
		element, remaining, err = parseDeserialize(remaining)
		if err != nil {
			return nil, remaining, err
		}
		value = append(value, element)
		length--
	}

	return &Data{kind: ArrayKind, value: value}, remaining, nil
}

func deserializeSimpleError(s string) (*Data, string, error) {
	m := SIMPLE_ERROR_RE.FindStringSubmatch(s)
	if m == nil {
		return nil, "", fmt.Errorf("error simple error format: %q", s)
	}
	value := m[1]
	remaining := m[2]
	return &Data{kind: SimpleErrorKind, value: value}, remaining, nil
}

func parseDeserialize(s string) (*Data, string, error) {
	firstChar := s[0]
	remaining := s[1:]
	switch firstChar {
	case '$':
		return deserializeBulkStringOrNull(remaining)
	case '+':
		return deserializeSimpleString(remaining)
	case ':':
		return deserializeInt(remaining)
	case '*':
		return deserializeArray(remaining)
	case '-':
		return deserializeSimpleError(remaining)
	default:
		return nil, "", fmt.Errorf("error unexpected first char: %q", firstChar)
	}
}

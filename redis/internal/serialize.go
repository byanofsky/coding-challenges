package internal

import (
	"fmt"
)

func Serialize(data Data) (string, error) {
	switch data.kind {
	case ArrayKind:
		return serializeArray(data)
	case NullKind:
		return serializeNull(), nil
	case SimpleStringKind:
		return serializeSimpleString(data)
	case BulkStringKind:
		return serializeBulkString(data)
	case IntKind:
		return serializeInt(data)
	case SimpleErrorKind:
		return serializeSimpleError(data)
	default:
		return "", fmt.Errorf("error unexpected data type: %v", data.kind)
	}
}

func serializeNull() string {
	return "$-1\r\n"
}

func serializeSimpleString(d Data) (string, error) {
	// TODO: Simple string validation
	s, err := d.GetString()
	if err != nil {
		return "", fmt.Errorf("error serializing simple string: %w", err)
	}
	// TODO: The string mustn't contain a CR (\r) or LF (\n) character and is terminated by CRLF (i.e., \r\n).
	return fmt.Sprintf("+%s\r\n", s), nil
}

func serializeBulkString(d Data) (string, error) {
	s, err := d.GetString()
	if err != nil {
		return "", fmt.Errorf("error serializing bulk string: %w", err)
	}
	l := len([]rune(s))
	return fmt.Sprintf("$%d\r\n%s\r\n", l, s), nil
}

func serializeInt(d Data) (string, error) {
	i, err := d.GetInt()
	if err != nil {
		return "", fmt.Errorf("error serializing int: %w", err)
	}
	return fmt.Sprintf(":%d\r\n", i), nil
}

func serializeArray(d Data) (string, error) {
	a, err := d.GetArray()
	if err != nil {
		return "", fmt.Errorf("error serializing array: %w", err)
	}

	result := fmt.Sprintf("*%d\r\n", len(*a))

	for _, item := range *a {
		s, err := Serialize(item)
		if err != nil {
			return "", err
		}
		result += s
	}

	return result, nil
}

func serializeSimpleError(d Data) (string, error) {
	s, err := d.GetString()
	if err != nil {
		return "", fmt.Errorf("error serializing simple string: %w", err)
	}
	return fmt.Sprintf("-%s\r\n", s), nil
}

package internal

import (
	"fmt"
)

type Kind int

const (
	NullKind Kind = iota
	StringKind
	IntKind
	ArrayKind
)

func (k Kind) String() string {
	switch k {
	case NullKind:
		return "Null"
	case StringKind:
		return "String"
	case IntKind:
		return "Int"
	case ArrayKind:
		return "Array"
	default:
		return "Unknown"
	}
}

type Data struct {
	kind  Kind
	value interface{}
}

func (d Data) String() string {
	switch d.kind {
	case NullKind:
		return "Data{Null}"
	case StringKind:
		return fmt.Sprintf("Data{%q}", d.value)
	case IntKind:
		return fmt.Sprintf("Data{%d}", d.value)
	case ArrayKind:
		return fmt.Sprintf("Data{%v}", d.value)
	default:
		return "Unknown"
	}
}

func (d Data) GetString() (string, error) {
	if d.kind != StringKind {
		return "", fmt.Errorf("cannot GetString of kind: %s", d.kind)
	}
	s, ok := d.value.(string)
	if !ok {
		return "", fmt.Errorf("error value is not a string: %v", d.value)
	}
	return s, nil
}

func (d Data) GetInt() (int, error) {
	if d.kind != IntKind {
		return 0, fmt.Errorf("cannot GetInt of kind: %s", d.kind)
	}
	i, ok := d.value.(int)
	if !ok {
		return 0, fmt.Errorf("error value is not an int: %v", d.value)
	}
	return i, nil
}

func (d Data) GetArray() (*[]Data, error) {
	if d.kind != ArrayKind {
		return nil, fmt.Errorf("cannot GetArray of kind: %s", d.kind)
	}
	a, ok := d.value.([]Data)
	if !ok {
		return nil, fmt.Errorf("error value is not an array: %v", d.value)
	}
	return &a, nil
}

func Serialize(data Data) (string, error) {
	switch data.kind {
	case ArrayKind:
		return serializeArray(data)
	case NullKind:
		return serializeNull(), nil
	case StringKind:
		return serializeString(data)
	case IntKind:
		return serializeInt(data)
	default:
		return "", fmt.Errorf("error unexpected data type: %v", data.kind)
	}
}

func serializeNull() string {
	return "$-1\r\n"
}

func serializeString(d Data) (string, error) {
	s, err := d.GetString()
	if err != nil {
		return "", fmt.Errorf("error serializing string: %w", err)
	}
	// TODO: The string mustn't contain a CR (\r) or LF (\n) character and is terminated by CRLF (i.e., \r\n).
	return fmt.Sprintf("+%s\r\n", s), nil
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

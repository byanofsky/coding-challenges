package internal

import (
	"fmt"
)

type Null struct{}

type Serializable interface {
	Null | int | string
}

func Serialize[T Serializable](data T) (string, error) {
	switch t := any(data).(type) {
	case Null:
		return serializeNull(), nil
	case string:
		s := any(data).(string)
		return serializeString(s), nil
	default:
		return "", fmt.Errorf("error unexpected data type: %v", t)
	}
}

func serializeNull() string {
	return "$-1\r\n"
}

func serializeString(s string) string {
	return fmt.Sprintf("+%s\r\n", s)
}

package internal

import (
	"fmt"
	"reflect"
)

type Null struct{}

type Serializable interface {
	Null | int | string | []any
}

func Serialize(data any) (string, error) {
	fmt.Println(reflect.TypeOf(data).Kind())
	if reflect.TypeOf(data).Kind() == reflect.Slice {
		a := any(data).([]any)
		return serializeArray(a)
	}

	switch t := any(data).(type) {
	case Null:
		return serializeNull(), nil
	case string:
		s := any(data).(string)
		return serializeString(s), nil
	case int:
		i := any(data).(int)
		return serializeInt(i), nil
	default:
		return "", fmt.Errorf("error unexpected data type: %v", t)
	}
}

func serializeNull() string {
	return "$-1\r\n"
}

func serializeString(s string) string {
	// TODO: The string mustn't contain a CR (\r) or LF (\n) character and is terminated by CRLF (i.e., \r\n).
	return fmt.Sprintf("+%s\r\n", s)
}

func serializeInt(i int) string {
	return fmt.Sprintf(":%d\r\n", i)
}

func serializeArray(a []any) (string, error) {
	result := fmt.Sprintf("*%d\r\n", len(a))

	for _, item := range a {
		s, err := Serialize(item)
		if err != nil {
			return "", err
		}
		result += s
	}

	return result, nil
}

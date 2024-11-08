package internal

import (
	"testing"
)

func TestSerialize(t *testing.T) {
	runSerializeTest(t, "Null", Data{kind: NullKind}, "$-1\r\n")
	runSerializeTest(t, "SimpleString1", Data{kind: StringKind, value: "OK"}, "+OK\r\n")
	runSerializeTest(t, "SimpleString2", Data{kind: StringKind, value: "hello world"}, "+hello world\r\n")
	runSerializeTest(t, "Int", Data{kind: IntKind, value: 5}, ":5\r\n")
	runSerializeTest(t, "Array SimpleString", Data{kind: ArrayKind, value: []Data{{kind: StringKind, value: "ping"}}}, "*1\r\n+ping\r\n")
	// runSerializeTest(t, "Array", []any{"ping"}, "*1\r\n$4\r\nping\r\n")
	// runSerializeTest(t, "", "", "*2\r\n$4\r\necho\r\n$11\r\nhello world\r\n")
	// runSerializeTest(t, "", "", "*2\r\n$3\r\nget\r\n$3\r\nkey\r\n")
	// runSerializeTest(t, "", "", "+OK\r\n")
	// runSerializeTest(t, "", "", "-Error message\r\n")
	// runSerializeTest(t, "", "", "$0\r\n\r\n")
}

func runSerializeTest(t *testing.T, name string, input Data, want string) {
	t.Run(name, func(t *testing.T) {
		got, err := Serialize(input)

		if err != nil {
			t.Fatalf("input %v, unexpected error: %v", input, err)
		}

		if got != want {
			t.Fatalf("input %v, want %q, got %q", input, want, got)
		}
	})
}

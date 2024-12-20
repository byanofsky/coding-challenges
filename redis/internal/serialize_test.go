package internal

import (
	"testing"
)

func TestSerialize(t *testing.T) {
	runSerializeTest(t, "Null", Data{kind: NullKind}, "$-1\r\n")
	runSerializeTest(t, "SimpleString1", Data{kind: SimpleStringKind, value: "OK"}, "+OK\r\n")
	runSerializeTest(t, "SimpleString2", Data{kind: SimpleStringKind, value: "hello world"}, "+hello world\r\n")
	// runSerializeTest(t, "SimpleStringInvalid1", Data{kind: StringKind, value: "hello\rworld"}, "+hello world\r\n")
	// runSerializeTest(t, "SimpleStringInvalid2", Data{kind: StringKind, value: "hello\nworld"}, "+hello world\r\n")
	runSerializeTest(t, "Int", Data{kind: IntKind, value: 5}, ":5\r\n")
	runSerializeTest(t, "Array SimpleString", Data{kind: ArrayKind, value: []Data{{kind: SimpleStringKind, value: "ping"}}}, "*1\r\n+ping\r\n")
	runSerializeTest(t, "BulkStringEmpty", Data{kind: BulkStringKind, value: ""}, "$0\r\n\r\n")
	runSerializeTest(t, "BulkString1", Data{kind: BulkStringKind, value: "hello world"}, "$11\r\nhello world\r\n")
	runSerializeTest(t, "Array BulkString1", Data{
		kind:  ArrayKind,
		value: []Data{{kind: BulkStringKind, value: "ping"}},
	}, "*1\r\n$4\r\nping\r\n")
	runSerializeTest(t, "Array BulkString2", Data{
		kind: ArrayKind,
		value: []Data{
			{kind: BulkStringKind, value: "echo"},
			{kind: BulkStringKind, value: "hello world"},
		},
	}, "*2\r\n$4\r\necho\r\n$11\r\nhello world\r\n")
	runSerializeTest(t, "Array BulkString3", Data{
		kind: ArrayKind,
		value: []Data{
			{kind: BulkStringKind, value: "get"},
			{kind: BulkStringKind, value: "key"},
		}}, "*2\r\n$3\r\nget\r\n$3\r\nkey\r\n")
	runSerializeTest(t, "Map",
		Data{
			kind: MapKind,
			value: map[Data]Data{
				*NewSimpleStringData("key1"): *NewBulkStringData("value1"),
			},
		}, "%1\r\n+key1\r\n$6\r\nvalue1\r\n")
	// TODO: map test with multiple elements. Need to handle ordering
	runSerializeTest(t, "SimpleError", Data{kind: SimpleErrorKind, value: "Error message"}, "-Error message\r\n")
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

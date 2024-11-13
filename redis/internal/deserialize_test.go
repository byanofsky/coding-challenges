package internal

import (
	"reflect"
	"testing"
)

func TestDeserialize(t *testing.T) {
	runDeserializeTest(t, "Null", "$-1\r\n", Data{kind: NullKind})
	runDeserializeTest(t, "SimpleString1", "+OK\r\n", Data{kind: SimpleStringKind, value: "OK"})
	runDeserializeTest(t, "SimpleString2", "+hello world\r\n", Data{kind: SimpleStringKind, value: "hello world"})
	// runSerializeTest(t, "SimpleStringInvalid1", Data{kind: StringKind, value: "hello\rworld"}, "+hello world\r\n")
	// runSerializeTest(t, "SimpleStringInvalid2", Data{kind: StringKind, value: "hello\nworld"}, "+hello world\r\n")
	runDeserializeTest(t, "Int", ":5\r\n", Data{kind: IntKind, value: 5})
	runDeserializeTest(t, "Int", ":+5\r\n", Data{kind: IntKind, value: 5})
	runDeserializeTest(t, "Int", ":-5\r\n", Data{kind: IntKind, value: -5})
	runDeserializeTest(t, "Array SimpleString", "*1\r\n+ping\r\n", Data{kind: ArrayKind, value: []Data{{kind: SimpleStringKind, value: "ping"}}})
	runDeserializeTest(t, "Array SimpleString", "*2\r\n+ping\r\n+pong\r\n", Data{
		kind:  ArrayKind,
		value: []Data{{kind: SimpleStringKind, value: "ping"}, {kind: SimpleStringKind, value: "pong"}}})
	runDeserializeTest(t, "BulkStringEmpty", "$0\r\n\r\n", Data{kind: BulkStringKind, value: ""})
	runDeserializeTest(t, "BulkString1", "$11\r\nhello world\r\n", Data{kind: BulkStringKind, value: "hello world"})
	runDeserializeTest(t, "BulkString With CRLF", "$12\r\nhello\r\nworld\r\n", Data{kind: BulkStringKind, value: "hello\r\nworld"})
	runDeserializeTest(t, "Array BulkString1", "*1\r\n$4\r\nping\r\n", Data{
		kind:  ArrayKind,
		value: []Data{{kind: BulkStringKind, value: "ping"}},
	})
	runDeserializeTest(t, "Array BulkString2", "*2\r\n$4\r\necho\r\n$11\r\nhello world\r\n", Data{
		kind: ArrayKind,
		value: []Data{
			{kind: BulkStringKind, value: "echo"},
			{kind: BulkStringKind, value: "hello world"},
		},
	})
	runDeserializeTest(t, "Array BulkString3", "*2\r\n$3\r\nget\r\n$3\r\nkey\r\n", Data{
		kind: ArrayKind,
		value: []Data{
			{kind: BulkStringKind, value: "get"},
			{kind: BulkStringKind, value: "key"},
		}})
	runDeserializeTest(t, "Map", "%2\r\n+key1\r\n$6\r\nvalue1\r\n+key2\r\n$6\r\nvalue2\r\n", Data{
		kind: MapKind,
		value: map[Data]Data{
			*NewSimpleStringData("key1"): *NewBulkStringData("value1"),
			*NewSimpleStringData("key2"): *NewBulkStringData("value2"),
		},
	})
	// runDeserializeTest(t, "SimpleError", "-Error message\r\n", &Data{kind: SimpleErrorKind, value: "Error message"})
}

func runDeserializeTest(t *testing.T, name string, input string, want Data) {
	t.Run(name, func(t *testing.T) {
		got, err := Deserialize(input)

		if err != nil {
			t.Fatalf("input %q, unexpected error: %v", input, err)
		}

		if !reflect.DeepEqual(*got, want) {
			t.Fatalf("input %q, want %v, got %v", input, want, got)
		}
	})
}

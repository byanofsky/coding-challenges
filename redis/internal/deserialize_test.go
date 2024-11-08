package internal

import (
	"reflect"
	"testing"
)

func TestDeserialize(t *testing.T) {
	runDeserializeTest(t, "Null", "$-1\r\n", Data{kind: NullKind})
	runDeserializeTest(t, "SimpleString1", "+OK\r\n", Data{kind: SimpleStringKind, value: "OK"})
	runDeserializeTest(t, "SimpleString2", "+hello world\r\n", Data{kind: SimpleStringKind, value: "hello world"})
	// // runSerializeTest(t, "SimpleStringInvalid1", Data{kind: StringKind, value: "hello\rworld"}, "+hello world\r\n")
	// // runSerializeTest(t, "SimpleStringInvalid2", Data{kind: StringKind, value: "hello\nworld"}, "+hello world\r\n")
	// runSerializeTest(t, "Int", Data{kind: IntKind, value: 5}, ":5\r\n")
	// runSerializeTest(t, "Array SimpleString", Data{kind: ArrayKind, value: []Data{{kind: SimpleStringKind, value: "ping"}}}, "*1\r\n+ping\r\n")
	// runSerializeTest(t, "BulkStringEmpty", Data{kind: BulkStringKind, value: ""}, "$0\r\n\r\n")
	// runSerializeTest(t, "BulkString1", Data{kind: BulkStringKind, value: "hello world"}, "$11\r\nhello world\r\n")
	// runSerializeTest(t, "Array BulkString1", Data{
	// 	kind:  ArrayKind,
	// 	value: []Data{{kind: BulkStringKind, value: "ping"}},
	// }, "*1\r\n$4\r\nping\r\n")
	// runSerializeTest(t, "Array BulkString2", Data{
	// 	kind: ArrayKind,
	// 	value: []Data{
	// 		{kind: BulkStringKind, value: "echo"},
	// 		{kind: BulkStringKind, value: "hello world"},
	// 	},
	// }, "*2\r\n$4\r\necho\r\n$11\r\nhello world\r\n")
	// runSerializeTest(t, "Array BulkString3", Data{
	// 	kind: ArrayKind,
	// 	value: []Data{
	// 		{kind: BulkStringKind, value: "get"},
	// 		{kind: BulkStringKind, value: "key"},
	// 	}}, "*2\r\n$3\r\nget\r\n$3\r\nkey\r\n")
	// runSerializeTest(t, "SimpleError", Data{kind: SimpleErrorKind, value: "Error message"}, "-Error message\r\n")
}

func runDeserializeTest(t *testing.T, name string, input string, want Data) {
	t.Run(name, func(t *testing.T) {
		got, err := Deserialize(input)

		if err != nil {
			t.Fatalf("input %v, unexpected error: %v", input, err)
		}

		if !reflect.DeepEqual(*got, want) {
			t.Fatalf("input %q, want %v, got %v", input, want, got)
		}
	})
}

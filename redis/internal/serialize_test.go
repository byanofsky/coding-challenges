package internal

import (
	"testing"
)

func TestSerialize(t *testing.T) {
	runSerializeTest(t, "Null", Null{}, "$-1\r\n")
	runSerializeTest(t, "String", "OK", "+OK\r\n")
	runSerializeTest(t, "Int", 5, ":5\r\n")
}

func runSerializeTest[T Serializable](t *testing.T, name string, input T, want string) {
	t.Run(name, func(t *testing.T) {
		got, err := Serialize(input)

		if err != nil {
			t.Fatalf("input %v, unexpected error: %v", input, err)
		}

		if got != want {
			t.Fatalf("input %v, want %v, got %v", input, want, got)
		}
	})
}

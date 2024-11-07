package internal

import (
	"testing"
)

func TestSerializeNull(t *testing.T) {
	input := Null{}
	want := "$-1\r\n"

	got, err := Serialize(input)

	if err != nil {
		t.Fatalf("input %v, unexpected error: %v", input, err)
	}

	if got != want {
		t.Fatalf("input %v, want %v, got %v", input, want, got)
	}
}

func TestSerializeString(t *testing.T) {
	input := "OK"
	want := "+OK\r\n"

	got, err := Serialize(input)

	if err != nil {
		t.Fatalf("input %v, unexpected error: %v", input, err)
	}

	if got != want {
		t.Fatalf("input %v, want %v, got %v", input, want, got)
	}
}

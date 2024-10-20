package my_regexp_test

import (
	my_regexp "my-regexp"
	"testing"
)

func TestCompile(t *testing.T) {
	expected := true
	actual := my_regexp.Compile()

	if expected != actual {
		t.Errorf("Expected: %v\nActual: %v", expected, actual)
	}
}

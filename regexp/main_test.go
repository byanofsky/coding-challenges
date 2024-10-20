package regexp_test

import (
	"testing"
)

func TestHelloWorld(t *testing.T) {
	expected := true
	actual := true

	if expected != actual {
		t.Errorf("Expected: %v\nActual: %v", expected, actual)
	}
}

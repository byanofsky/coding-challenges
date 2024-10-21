package my_regexp

import "testing"

func TestCompileMatchSuccess(t *testing.T) {
	re := Compile("abc")
	actual, err := re.Match("abc")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := true

	if expected != actual {
		t.Errorf("Expected: %v\nActual: %v", expected, actual)
	}
}

func TestCompileMatchFail(t *testing.T) {
	re := Compile("abc")
	actual, err := re.Match("def")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := false

	if expected != actual {
		t.Errorf("Expected: %v\nActual: %v", expected, actual)
	}
}

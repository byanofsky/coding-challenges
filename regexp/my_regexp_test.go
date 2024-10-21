package my_regexp

import "testing"

func TestCompileMatchSuccess(t *testing.T) {
	re, err := Compile("abc")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
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
	re, err := Compile("abc")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	actual, err := re.Match("def")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := false

	if expected != actual {
		t.Errorf("Expected: %v\nActual: %v", expected, actual)
	}
}

func TestWildcard(t *testing.T) {
	re, err := Compile(".b")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	actual, err := re.Match("ab")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := true

	if expected != actual {
		t.Errorf("Expected: %v\nActual: %v", expected, actual)
	}
}

// Tests that patter after wildcard must match
func TestWildcardNoMatch(t *testing.T) {
	re, err := Compile(".c")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	actual, err := re.Match("ab")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := false

	if expected != actual {
		t.Errorf("Expected: %v\nActual: %v", expected, actual)
	}
}

func TestRepetition(t *testing.T) {
	re, err := Compile("b*")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		t.FailNow()
	}
	actual, err := re.Match("bbb")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := true

	if expected != actual {
		t.Errorf("Expected: %v\nActual: %v", expected, actual)
	}
}

func TestRepetitionZeroTimesMatch(t *testing.T) {
	re, err := Compile("b*")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		t.FailNow()
	}
	actual, err := re.Match("aaa")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := true

	if expected != actual {
		t.Errorf("Expected: %v\nActual: %v", expected, actual)
	}
}

// Test when repition pattern after the repition
func TestRepetitionNoMatchPatternAfter(t *testing.T) {
	re, err := Compile("b*c")
	if err != nil {
		// TODO: Update tests to fail now
		t.Errorf("Unexpected error: %v", err)
		t.FailNow()
	}
	actual, err := re.Match("bbb")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := false

	if expected != actual {
		t.Errorf("Expected: %v\nActual: %v", expected, actual)
	}
}

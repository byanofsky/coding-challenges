package my_regexp

import (
	"reflect"
	"testing"
)

func TestScan(t *testing.T) {
	expected := []token{
		{kind: SingleCharacter, token: "a"},
		{kind: SingleCharacter, token: "."},
		{kind: Repetition, token: "*"},
		{kind: EscapeSequence, token: `\w`},
		{kind: EscapeSequence, token: `\n`},
	}
	actual := scan(`a.*\w\n`)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %v\nActual: %v", expected, actual)
	}
}

func TestSuccessSingleCharacterMatcher(t *testing.T) {
	m := newSingleCharacterMatcher('b')

	// TODO: Replace with mocking lib
	mockCalled := false
	var mockInt int
	mock := func(s string, i int) (bool, error) {
		mockCalled = true
		mockInt = i
		return true, nil
	}

	actual, err := m.isMatch("abc", 1, mock)
	expected := true

	if actual != expected {
		t.Errorf("Expected: %v\nActual: %v", expected, actual)
	}

	if err != nil {
		t.Errorf("Received error: %v", err)
	}

	if mockCalled != true {
		t.Errorf("Next function not called")
	}

	if mockInt != 2 {
		t.Errorf("Next function called with wrong arg. Expected: 2. Actual: %d", mockInt)
	}
}

func TestNoMatchSingleCharacterMatcher(t *testing.T) {
	m := newSingleCharacterMatcher('b')

	// TODO: Replace with mocking lib
	mockCalled := false
	mock := func(s string, i int) (bool, error) {
		mockCalled = true
		return true, nil
	}

	actual, err := m.isMatch("abc", 0, mock)
	expected := false

	if actual != expected {
		t.Errorf("Expected: %v\nActual: %v", expected, actual)
	}

	if err != nil {
		t.Errorf("Received error: %v", err)
	}

	if mockCalled != false {
		t.Errorf("Next function not called")
	}
}

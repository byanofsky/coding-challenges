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
	var mockMIdx int
	mock := func(s string, i int, mIdx int) (bool, error) {
		mockCalled = true
		mockInt = i
		mockMIdx = mIdx
		return true, nil
	}

	actual, err := m.isMatch("abc", 1, 0, mock)
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

	if mockMIdx != 1 {
		t.Errorf("Next function called with wrong arg. Expected: 1. Actual: %d", mockInt)
	}
}

func TestNoMatchSingleCharacterMatcher(t *testing.T) {
	m := newSingleCharacterMatcher('b')

	// TODO: Replace with mocking lib
	mockCalled := false
	mock := func(s string, i int, mIdx int) (bool, error) {
		mockCalled = true
		return true, nil
	}

	actual, err := m.isMatch("abc", 0, 0, mock)
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

func TestSuccessExactMatch(t *testing.T) {
	matchers := []matcher{
		newSingleCharacterMatcher('a'),
		newSingleCharacterMatcher('b'),
		newSingleCharacterMatcher('c'),
	}

	actual, err := isMatch(matchers, "abc")
	expected := true

	if err != nil {
		t.Errorf("Received error: %v", err)
	}

	if actual != expected {
		t.Errorf("Expected: %v\nActual: %v", expected, actual)
	}
}

func TestSuccessSubMatch(t *testing.T) {
	matchers := []matcher{
		newSingleCharacterMatcher('a'),
		newSingleCharacterMatcher('b'),
		newSingleCharacterMatcher('c'),
	}

	actual, err := isMatch(matchers, "123abc456")
	expected := true

	if err != nil {
		t.Errorf("Received error: %v", err)
	}

	if actual != expected {
		t.Errorf("Expected: %v\nActual: %v", expected, actual)
	}
}

func TestSuccessStartMatch(t *testing.T) {
	matchers := []matcher{
		newSingleCharacterMatcher('a'),
		newSingleCharacterMatcher('b'),
		newSingleCharacterMatcher('c'),
	}

	actual, err := isMatch(matchers, "123abc")
	expected := true

	if err != nil {
		t.Errorf("Received error: %v", err)
	}

	if actual != expected {
		t.Errorf("Expected: %v\nActual: %v", expected, actual)
	}
}

func TestSuccessEndMatch(t *testing.T) {
	matchers := []matcher{
		newSingleCharacterMatcher('a'),
		newSingleCharacterMatcher('b'),
		newSingleCharacterMatcher('c'),
	}

	actual, err := isMatch(matchers, "abc456")
	expected := true

	if err != nil {
		t.Errorf("Received error: %v", err)
	}

	if actual != expected {
		t.Errorf("Expected: %v\nActual: %v", expected, actual)
	}
}

func TestNoMatchMultipleCharactersMatch(t *testing.T) {
	matchers := []matcher{
		newSingleCharacterMatcher('a'),
		newSingleCharacterMatcher('b'),
		newSingleCharacterMatcher('c'),
	}

	actual, err := isMatch(matchers, "abd")
	expected := false

	if err != nil {
		t.Errorf("Received error: %v", err)
	}

	if actual != expected {
		t.Errorf("Expected: %v\nActual: %v", expected, actual)
	}
}

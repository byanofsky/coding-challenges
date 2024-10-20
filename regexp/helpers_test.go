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

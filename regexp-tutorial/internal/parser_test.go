package internal

import (
	"reflect"
	"testing"
)

func TestParseRangeQuantifier(t *testing.T) {
	tests := []struct {
		name  string
		input string
		match bool
		rq    *RangeQuantifier
		after string
	}{{
		name:  "positive: basic case",
		input: "{3}",
		match: true,
		rq:    &RangeQuantifier{lowerBound: 3, upperBound: 0},
		after: "",
	}, {
		name:  "negative: no braces case",
		input: "3",
		match: false,
		rq:    nil,
		after: "3",
	}, {
		name:  "negative: no number",
		input: "{abc}",
		match: false,
		rq:    nil,
		after: "{abc}",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, rq, after := parseRangeQuantifier(tt.input)
			if match != tt.match {
				t.Fatalf("input %q, match %v, want %v", tt.input, match, tt.match)
			}
			if !reflect.DeepEqual(rq, tt.rq) {
				t.Errorf("input %q, rq %v, want %v", tt.input, *rq, *tt.rq)
			}
			if after != tt.after {
				t.Errorf("input %q, after %q, want %q", tt.input, after, tt.after)
			}
		})

	}
}

func TestParseCharacter(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		result rune
		found  bool
		after  string
	}{{
		name:   "positive: basic case",
		input:  "abc",
		result: 'a',
		found:  true,
		after:  "bc",
	}, {
		name:   "negative: empty",
		input:  "",
		result: 0,
		found:  false,
		after:  "",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewCharacterParser()
			result, after, found, err := p.parse(tt.input)
			if err != nil {
				t.Fatalf("input %q unexpected error: %v", tt.input, err)
			}
			if result != tt.result {
				t.Fatalf("input %q, result %v, want %v", tt.input, result, tt.result)
			}
			if found != tt.found {
				t.Fatalf("input %q, found %v, want %v", tt.input, found, tt.found)
			}
			if after != tt.after {
				t.Errorf("input %q, after %q, want %q", tt.input, after, tt.after)
			}
		})
	}
}

func TestParseZeroOrMoreCharacter(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		result []rune
		found  bool
		after  string
	}{{
		name:   "positive: basic case",
		input:  "abc",
		result: []rune{'a', 'b', 'c'},
		found:  true,
		after:  "",
	}, {
		name:   "positive: zero result match",
		input:  "",
		result: []rune{},
		found:  true,
		after:  "",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewZeroOrMoreParser(NewCharacterParser())
			result, after, found, err := p.parse(tt.input)
			if err != nil {
				t.Fatalf("input %q unexpected error: %v", tt.input, err)
			}
			if !reflect.DeepEqual(result, tt.result) {
				t.Fatalf("input %q, result %v, want %v", tt.input, result, tt.result)
			}
			if found != tt.found {
				t.Fatalf("input %q, found %v, want %v", tt.input, found, tt.found)
			}
			if after != tt.after {
				t.Errorf("input %q, after %q, want %q", tt.input, after, tt.after)
			}
		})
	}
}

func TestParseOneOrMoreCharacter(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		result []rune
		found  bool
		after  string
	}{{
		name:   "positive: basic case",
		input:  "abc",
		result: []rune{'a', 'b', 'c'},
		found:  true,
		after:  "",
	}, {
		name:   "positive: zero result match",
		input:  "",
		result: nil,
		found:  false,
		after:  "",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewOneOrMoreParser(NewCharacterParser())
			result, after, found, err := p.parse(tt.input)
			if err != nil {
				t.Fatalf("input %q unexpected error: %v", tt.input, err)
			}
			if !reflect.DeepEqual(result, tt.result) {
				t.Fatalf("input %q, result %v, want %v", tt.input, result, tt.result)
			}
			if found != tt.found {
				t.Fatalf("input %q, found %v, want %v", tt.input, found, tt.found)
			}
			if after != tt.after {
				t.Errorf("input %q, after %q, want %q", tt.input, after, tt.after)
			}
		})
	}
}

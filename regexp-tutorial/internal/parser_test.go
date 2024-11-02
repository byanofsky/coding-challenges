package internal

import (
	"reflect"
	"testing"
)

func TestParseRangeQuantifier(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		result RangeQuantifier
		found  bool
		after  string
	}{{
		name:   "positive: basic case",
		input:  "{3,5}",
		result: RangeQuantifier{LowerBound: 3, UpperBound: 5},
		found:  true,
		after:  "",
	}, {
		name:   "negative: no braces case",
		input:  "3",
		result: RangeQuantifier{LowerBound: 0, UpperBound: 0},
		found:  false,
		after:  "",
	}, {
		name:   "negative: no number",
		input:  "{abc}",
		result: RangeQuantifier{LowerBound: 0, UpperBound: 0},
		found:  false,
		after:  "",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, after, found, err := NewRangeQuantifier().parse(tt.input)
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

func TestDigitParser(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		result rune
		found  bool
		after  string
	}{{
		name:   "positive: basic case",
		input:  "1bc",
		result: '1',
		found:  true,
		after:  "bc",
	}, {
		name:   "negative: not digit",
		input:  "abc",
		result: 0,
		found:  false,
		// TODO: Should return full string when not found
		after: "bc",
	}, {
		name:   "negative: empty",
		input:  "",
		result: 0,
		found:  false,
		after:  "",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewDigitParser()
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

func TestNumberParser(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		result int
		found  bool
		after  string
	}{{
		name:   "positive: basic case",
		input:  "123",
		result: 123,
		found:  true,
		after:  "",
	}, {
		name:   "positive: with after",
		input:  "123abc",
		result: 123,
		found:  true,
		after:  "abc",
	}, {
		name:   "negative: not number",
		input:  "abc",
		result: 0,
		found:  false,
		// TODO: Should return full string when not found
		after: "",
	}, {
		name:   "negative: empty",
		input:  "",
		result: 0,
		found:  false,
		after:  "",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewNumberParser()
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

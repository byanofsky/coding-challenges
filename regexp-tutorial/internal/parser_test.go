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
		result: RangeQuantifier{LowerBound: 3, UpperBound: NewOptionalVal(5)},
		found:  true,
		after:  "",
	}, {
		name:   "positive: optional",
		input:  "{3,}",
		result: RangeQuantifier{LowerBound: 3, UpperBound: NewEmptyOptionalVal[int]()},
		found:  true,
		after:  "",
	}, {
		name:   "positive: exact match",
		input:  "{3}",
		result: RangeQuantifier{LowerBound: 3, UpperBound: NewOptionalVal(3)},
		found:  true,
		after:  "",
	}, {
		name:   "negative: not range quantifier",
		input:  "3",
		result: RangeQuantifier{LowerBound: 0, UpperBound: OptionalVal[int]{}},
		found:  false,
		after:  "",
	}, {
		name:   "negative: no number",
		input:  "{abc}",
		result: RangeQuantifier{LowerBound: 0, UpperBound: OptionalVal[int]{}},
		found:  false,
		after:  "",
	}, {
		name:   "negative: no brace",
		input:  "{3,",
		result: RangeQuantifier{LowerBound: 0, UpperBound: OptionalVal[int]{}},
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

func TestOrThrowParser(t *testing.T) {
	t.Run("test OrThrow parser panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected function to panic, but it did not")
			} else if r != "expected number" {
				t.Errorf("Expect 'expected number', got '%v'", r)
			}
		}()

		p := OrThrow(NewNumberParser(), "expected number")
		p.parse("abc")
	})
}

func TestOneOf(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		result MatchResult
		found  bool
		after  string
	}{{
		name:   "positive: match first parser",
		input:  "abc",
		result: MatchResult{kind: MatchResultString, s: "a"},
		found:  true,
		after:  "bc",
	}, {
		name:   "positive: match second parser",
		input:  "123",
		result: MatchResult{kind: MatchResultNumber, n: 123},
		found:  true,
		after:  "",
	}, {
		name:   "negative: no match",
		input:  "()",
		result: MatchResult{},
		found:  false,
		after:  "()",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := OneOf(Wrap(NewStringParser("a")), Wrap(NewNumberParser()))
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

func TestExpression(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		result Expression
		found  bool
		after  string
	}{{
		name:  "positive:",
		input: "ab12",
		result: Expression{matches: []MatchResult{
			{kind: MatchResultRune, r: 'a'},
			{kind: MatchResultRune, r: 'b'},
			{kind: MatchResultNumber, n: 12},
		}},
		found: true,
		after: "",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewExpression()
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

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

func TestParseRangeQuantifierNegative(t *testing.T) {
	input := "3"
	match, _, after := parseRangeQuantifier(input)

	wantMatch := false
	wantAfter := input

	if match != wantMatch {
		t.Errorf("input %q, match %v, want %v", input, match, wantMatch)
	}
	if after != wantAfter {
		t.Errorf("input %q, after %v, want %v", input, after, wantAfter)
	}
}

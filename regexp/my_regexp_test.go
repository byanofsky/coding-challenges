package my_regexp

import "testing"

func TestCompileMatch(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		input   string
		want    bool
	}{
		{
			name:    "exact match",
			pattern: "abc",
			input:   "abc",
			want:    true,
		},
		{
			name:    "no match",
			pattern: "abc",
			input:   "def",
			want:    false,
		},
		{
			name:    "wildcard match",
			pattern: "a.c",
			input:   "abc",
			want:    true,
		},
		{
			// TODO: Add similar test when pattern not found in substring
			name:    "wildcard no match",
			pattern: "a.c",
			input:   "def",
			want:    false,
		},
		{
			name:    "star repitition match",
			pattern: "b*",
			input:   "bbb",
			want:    true,
		},
		{
			name:    "star repitition zero match",
			pattern: "b*",
			input:   "aaa",
			want:    true,
		},
		{
			name:    "star repitition no match",
			pattern: "b*c",
			input:   "aaa",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := mustCompile(t, tt.pattern)

			got, err := re.Match(tt.input)
			if err != nil {
				t.Fatalf("Match(%q) unexpected error: %v", tt.input, err)
			}

			if got != tt.want {
				t.Errorf("Match(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func mustCompile(t *testing.T, pattern string) *Regexp {
	t.Helper()

	re, err := Compile(pattern)
	if err != nil {
		t.Fatalf("Compile(%q) unexpected error: %v", pattern, err)
	}
	return re
}

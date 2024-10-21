package my_regexp

import (
	"fmt"
)

// Given a regexp pattern, returns slice containing tokens.
// This function is used during compilation.
func scan(pattern string) []token {
	tokens := make([]token, 0)

	isEscapeSeq := false
	for _, char := range pattern {
		var t string
		var kind tokenType

		if isEscapeSeq {
			isEscapeSeq = false
			kind = EscapeSequence
			t = fmt.Sprintf(`\%c`, char)
		} else {
			switch char {
			case '\\':
				isEscapeSeq = true
				continue
			case '*':
				t = "*"
				kind = Repetition
			default:
				t = string(char)
				kind = SingleCharacter
			}
		}

		tokens = append(tokens, token{kind: kind, token: t})
	}
	return tokens
}

func newSingleCharacterMatcher(c byte) matcher {
	f := func(s string, i int, next nextMatcher) (bool, error) {
		sLen := len(s)

		// Assert edge cases
		if i >= sLen {
			// TODO: Custom error
			return false, fmt.Errorf("out of bounds: %d of %d", i, sLen)
		}

		if s[i] != c {
			return false, nil
		}

		return next(s, i+1)
	}
	return matcher{isMatch: f}
}

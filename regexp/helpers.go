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

func isMatch(matchers []matcher, s string) (bool, error) {
	var next nextMatcher
	next = func(s string, i int, mIdx int) (bool, error) {
		// Reached end of pattern. Therefore, this is a match.
		if mIdx == len(matchers) {
			return true, nil
		}

		return matchers[mIdx].isMatch(s, i, mIdx, next)
	}

	for i := range s {
		// TODO: Pass pointer to matchers
		result, err := matchers[0].isMatch(s, i, 0, next)
		if err != nil {
			return false, err
		}
		if result {
			return true, nil
		}
	}

	return false, nil
}

func newSingleCharacterMatcher(c byte) matcher {
	f := func(s string, i int, mIdx int, next nextMatcher) (bool, error) {
		sLen := len(s)

		// Assert edge cases
		if i >= sLen {
			// TODO: Custom error
			return false, fmt.Errorf("out of bounds: %d of %d", i, sLen)
		}

		if s[i] != c {
			return false, nil
		}

		return next(s, i+1, mIdx+1)
	}
	return matcher{isMatch: f}
}

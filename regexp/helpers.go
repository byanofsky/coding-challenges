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
			case '+':
				t = "+"
				kind = Repetition
			case '.':
				t = "."
				kind = Wildcard
			default:
				t = string(char)
				kind = SingleCharacter
			}
		}

		tokens = append(tokens, token{kind: kind, token: t})
	}
	return tokens
}

type repitionConstructor = func(matcher) matcher

// TODO: Return pointer to matcher
func parse(tokens []token) ([]matcher, error) {
	// TODO: Move to top level definition
	repitionMap := map[string]repitionConstructor{
		"*": newStarRepition,
		"+": newPlusRepitition,
	}

	matchers := make([]matcher, 0)

	i := 0
	for i < len(tokens) {
		token := tokens[i]
		var m matcher

		switch token.kind {
		case SingleCharacter:
			m = newSingleCharacterMatcher(token.token[0])
		case Wildcard:
			m = newWildcardMatcher()
		case Repetition:
			return nil, fmt.Errorf("possible reptition followed by repition. index: %d", i)
		default:
			return nil, fmt.Errorf("unknown pattern: %s. index: %d", token.token, i)
		}

		if i+1 < len(tokens) {
			// Peak next char to check for repition
			nextToken := tokens[i+1].token
			if rm, ok := repitionMap[nextToken]; ok {
				matchers = append(matchers, rm(m))
				// Skip to char after repition
				i = i + 2
				continue
			}
		}

		matchers = append(matchers, m)
		i++
	}
	return matchers, nil
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

package my_regexp

import "fmt"

func (r Regexp) Match(s string) (bool, error) {
	return isMatch(r.matchers, s)
}

func Compile(pattern string) (*Regexp, error) {
	tokens := scan(pattern)

	matchers := make([]matcher, 0)

	i := 0
	for i < len(tokens) {
		token := tokens[i]
		var matcher matcher

		switch token.kind {
		case SingleCharacter:
			matcher = newSingleCharacterMatcher(token.token[0])
		case Wildcard:
			matcher = newWildcardMatcher()
		case Repetition:
			return nil, fmt.Errorf("possible reptition followed by repition. index: %d", i)
		default:
			return nil, fmt.Errorf("unknown pattern: %s. index: %d", token.token, i)
		}

		if i+1 < len(tokens) {
			// Peak next char to check for repition
			// TODO: Abstrac to is Reptition function
			if tokens[i+1].token == string('*') {
				matcher = newStarRepition(matcher)
				i++
			}
		}

		matchers = append(matchers, matcher)
		i++
	}

	re := Regexp{matchers: matchers}
	return &re, nil
}

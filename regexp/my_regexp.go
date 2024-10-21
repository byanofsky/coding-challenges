package my_regexp

func (r Regexp) Match(s string) (bool, error) {
	return isMatch(r.matchers, s)
}

func Compile(pattern string) *Regexp {
	tokens := scan(pattern)

	matchers := make([]matcher, 0)
	for _, token := range tokens {
		switch token.kind {
		case SingleCharacter:
			matchers = append(matchers, newSingleCharacterMatcher(token.token[0]))
		case Wildcard:
			matchers = append(matchers, newWildcardMatcher())
		}
	}

	re := Regexp{matchers: matchers}
	return &re
}

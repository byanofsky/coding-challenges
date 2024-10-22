package my_regexp

func (r Regexp) Match(s string) (bool, error) {
	return isMatch(r.matchers, s)
}

func Compile(pattern string) (*Regexp, error) {
	// TODO: Scan should return err
	tokens := scan(pattern)
	matchers, err := parse(tokens)
	if err != nil {
		return nil, err
	}

	re := Regexp{matchers: matchers}
	return &re, nil
}

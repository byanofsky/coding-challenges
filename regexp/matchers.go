package my_regexp

func newSingleCharacterMatcher(c byte) matcher {
	f := func(s string, i int, mIdx int, next nextMatcher) (bool, error) {
		// TODO: Move into next matcher function.
		// Although, such an error may not be needed. This may be a valid case to return false.
		// Assert edge cases
		if i >= len(s) {
			// TODO: Custom error
			return false, nil
		}

		if s[i] != c {
			return false, nil
		}

		return next(s, i+1, mIdx+1)
	}
	return matcher{isMatch: f}
}

func newWildcardMatcher() matcher {
	f := func(s string, i int, mIdx int, next nextMatcher) (bool, error) {
		if i >= len(s) {
			return false, nil
		}
		return next(s, i+1, mIdx+1)
	}

	return matcher{isMatch: f}
}

func newStarRepition(m matcher) matcher {
	f := func(s string, i int, mIdx int, next nextMatcher) (bool, error) {
		subi := i
		for subi < len(s) {
			// The next function always returns true, because next only called if matcher matched
			result, err := m.isMatch(s, subi, mIdx, func(s string, i, mIdx int) (bool, error) { return true, nil })
			if err != nil {
				return false, err
			}
			if !result {
				break
			}
			subi++
		}

		// subi will now be either out of bounds of list or index that didn't match

		for subi >= i {
			result, err := next(s, subi, mIdx+1)
			if err != nil {
				return false, err
			}
			if result {
				return result, nil
			}
			subi--
		}

		return false, nil
	}

	return matcher{isMatch: f}
}
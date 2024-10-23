package internal

import (
	"strconv"
	"strings"
	"unicode"
)

type RangeQuantifier struct {
	lowerBound int
	upperBound int
}

func parseString(prefix string, s string) (after string, found bool) {
	return strings.CutPrefix(s, prefix)
}

func parseNumber(s string) (result int, after string, found bool) {
	var i int
	var c rune
	for i, c = range s {
		if !unicode.IsDigit(c) {
			break
		}
	}
	after = s[i:]

	result, err := strconv.Atoi(s[0:i])
	if err != nil {
		return 0, after, false
	}
	return result, after, true
}

/*
Returns a tuple of three values:
 1. A boolean indicating whether a match was found
 2. A RangeQuantifier
 3. The remaining string

2 and 3 are returned iff a match is found.
*/
func parseRangeQuantifier(s string) (match bool, rq *RangeQuantifier, after string) {
	after, found := parseString("{", s)
	if !found {
		return false, nil, s
	}

	lowerBound, after, found := parseNumber(after)
	if !found {
		return false, nil, s
	}

	after, found = parseString("}", after)
	if !found {
		return false, nil, s
	}

	return true, &RangeQuantifier{lowerBound: lowerBound, upperBound: 0}, after
}

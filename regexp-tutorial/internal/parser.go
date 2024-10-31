package internal

import (
	"strconv"
	"strings"
	"unicode"
)

type Parser[A any] struct {
	parse func(s string) (result A, substring string, found bool, err error)
}

// func (p Parser[A]) filter(predicate func()) Parser[A] {
// 	return Parser[A]{}
// }

func FlatMap[A any, B any](p Parser[A], transform func(A) Parser[B]) Parser[B] {
	return Parser[B]{
		parse: func(s string) (result B, substring string, found bool, err error) {
			a, intermediate, pFound, _ := p.parse(s)
			if !pFound {
				return result, substring, false, nil
			}
			return transform(a).parse(intermediate)
		},
	}
}

func Map[A any, B any](p Parser[A], transform func(A) (b B, found bool)) Parser[B] {
	return FlatMap(p, func(match A) Parser[B] {
		return Parser[B]{
			parse: func(s string) (result B, substring string, found bool, err error) {
				result, found = transform(match)
				return result, s, found, nil
			},
		}
	})
}

func NewOneOrMoreParser[A any](p Parser[A]) Parser[[]A] {
	zeroOrMore := NewZeroOrMoreParser(p)
	return Map(zeroOrMore, func(match []A) ([]A, bool) {
		if len(match) == 0 {
			// Allow return nil
			return nil, false
		}
		return match, true
	})
}

func NewZeroOrMoreParser[A any](p Parser[A]) Parser[[]A] {
	return Parser[[]A]{
		parse: func(s string) (result []A, substring string, found bool, err error) {
			rest := s
			matches := make([]A, 0, len(s))
			for {
				r, ss, found, _ := p.parse(rest)
				if !found {
					break
				}
				matches = append(matches, r)
				rest = ss
			}
			return matches, rest, true, nil
		},
	}
}

func NewStringParser(p string) Parser[string] {
	return Parser[string]{
		parse: func(s string) (result string, substring string, found bool, err error) {
			result = ""
			substring, found = strings.CutPrefix(s, p)
			return
		},
	}

}

func NewCharacterParser() Parser[rune] {
	return Parser[rune]{
		parse: func(s string) (result rune, substring string, found bool, err error) {
			if len(s) == 0 {
				return
			}
			found = true
			runes := []rune(s)
			result = runes[0]
			substring = string(runes[1:])
			return
		},
	}
}

type RangeQuantifier struct {
	lowerBound int
	upperBound int
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
	_, after, found, _ := NewStringParser("{").parse(s)
	if !found {
		return false, nil, s
	}

	lowerBound, after, found := parseNumber(after)
	if !found {
		return false, nil, s
	}

	_, after, found, _ = NewStringParser("}").parse(after)
	if !found {
		return false, nil, s
	}

	return true, &RangeQuantifier{lowerBound: lowerBound, upperBound: 0}, after
}

package internal

import (
	"strconv"
	"strings"
	"unicode"
)

type Parser[A any] struct {
	parse func(s string) (result A, substring string, found bool, err error)
}

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

// TODO: Remove found return value
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

func Filter[A any](p Parser[A], pred func(A) bool) Parser[A] {
	return Map(p, func(a A) (b A, found bool) {
		if pred(a) {
			b = a
			found = true
		}
		return b, found
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

func NewDigitParser() Parser[rune] {
	p := NewCharacterParser()
	return Filter(p, func(a rune) bool {
		return unicode.IsDigit(a)
	})
}

func NewNumberParser() Parser[int] {
	p := NewDigitParser()
	q := NewOneOrMoreParser(p)
	return Map(q, func(a []rune) (b int, found bool) {
		s := string(a)
		n, err := strconv.Atoi(s)
		if err != nil {
			return n, false
		}
		return n, true
	})
}

type Zipped[A any, B any] struct {
	a A
	b B
}

func Zip[A any, B any](ap Parser[A], bp Parser[B]) Parser[Zipped[A, B]] {
	return FlatMap(ap, func(matchA A) Parser[Zipped[A, B]] {
		return Map(bp, func(matchB B) (Zipped[A, B], bool) {
			return Zipped[A, B]{matchA, matchB}, true
		})
	})
}

type Zipped3[A any, B any, C any] struct {
	a A
	b B
	c C
}

func Zip3[A any, B any, C any](ap Parser[A], bp Parser[B], cp Parser[C]) Parser[Zipped3[A, B, C]] {
	return Map(Zip(ap, Zip(bp, cp)), func(z Zipped[A, Zipped[B, C]]) (Zipped3[A, B, C], bool) {
		return Zipped3[A, B, C]{
			z.a,
			z.b.a,
			z.b.b,
		}, true
	})
}

type Zipped4[A any, B any, C any, D any] struct {
	a A
	b B
	c C
	d D
}

func Zip4[A any, B any, C any, D any](ap Parser[A], bp Parser[B], cp Parser[C], dp Parser[D]) Parser[Zipped4[A, B, C, D]] {
	return Map(Zip3(ap, bp, Zip(cp, dp)), func(z Zipped3[A, B, Zipped[C, D]]) (Zipped4[A, B, C, D], bool) {
		return Zipped4[A, B, C, D]{
			z.a,
			z.b,
			z.c.a,
			z.c.b,
		}, true
	})
}

type Zipped5[A any, B any, C any, D any, E any] struct {
	a A
	b B
	c C
	d D
	e E
}

func Zip5[A any, B any, C any, D any, E any](ap Parser[A], bp Parser[B], cp Parser[C], dp Parser[D], ep Parser[E]) Parser[Zipped5[A, B, C, D, E]] {
	return Map(Zip4(ap, bp, cp, Zip(dp, ep)), func(z Zipped4[A, B, C, Zipped[D, E]]) (Zipped5[A, B, C, D, E], bool) {
		return Zipped5[A, B, C, D, E]{
			z.a,
			z.b,
			z.c,
			z.d.a,
			z.d.b,
		}, true
	})
}

type OptionalVal[A any] struct {
	none  bool
	value A
}

func Optional[A any](p Parser[A]) Parser[OptionalVal[A]] {
	return Parser[OptionalVal[A]]{
		parse: func(s string) (result OptionalVal[A], substring string, found bool, err error) {
			r, substring, f, _ := p.parse(s)
			if !f {
				return OptionalVal[A]{none: true}, s, true, nil
			}
			return OptionalVal[A]{none: false, value: r}, substring, true, nil
		},
	}
}

type RangeQuantifier struct {
	LowerBound int
	HasUpper   bool
	UpperBound int
}

/*
Returns a tuple of three values:
 1. A boolean indicating whether a match was found
 2. A RangeQuantifier
 3. The remaining string

2 and 3 are returned iff a match is found.
*/
func NewRangeQuantifier() Parser[RangeQuantifier] {
	o := Optional(Zip(NewStringParser(","), Optional(NewNumberParser())))
	p := Zip4(NewStringParser("{"), NewNumberParser(), o, NewStringParser("}"))
	return Map(p, func(z Zipped4[string, int, OptionalVal[Zipped[string, OptionalVal[int]]], string]) (RangeQuantifier, bool) {
		lowerBound := z.b
		upperOptional := z.c

		upperBound := 0
		hasUpper := false

		// No upper bound means exactly n times
		if upperOptional.none {
			hasUpper = true
			upperBound = lowerBound
		} else {
			if upperOptional.value.b.none {
				hasUpper = false
			} else {
				hasUpper = true
				upperBound = upperOptional.value.b.value
			}
		}

		return RangeQuantifier{LowerBound: lowerBound, UpperBound: upperBound, HasUpper: hasUpper}, true
	})
}

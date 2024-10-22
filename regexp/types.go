package my_regexp

type tokenType string

const (
	SingleCharacter tokenType = "SingleCharacter"
	Wildcard        tokenType = "Wildcard"
	EscapeSequence  tokenType = "EscapeSequence"
	Repetition      tokenType = "Repitition"
)

type token struct {
	kind  tokenType
	token string
}

type matcher struct {
	isMatch func(s string, i int, mIdx int, next nextMatcher) (bool, error)
}

type nextMatcher = func(s string, i int, mIdx int) (bool, error)

type Regexp struct {
	// TODO: Destroy during re destruction
	matchers []matcher
}

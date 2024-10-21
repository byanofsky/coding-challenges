package my_regexp

type tokenType int

const (
	SingleCharacter tokenType = iota
	Wildcard
	EscapeSequence
	Repetition
)

func (t tokenType) String() string {
	return [...]string{"SingleCharacter", "Wildcard", "EscapeSequence", "Repetition"}[t]
}

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

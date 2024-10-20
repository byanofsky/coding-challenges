package my_regexp

type tokenType int

const (
	SingleCharacter tokenType = iota
	EscapeSequence
	Repetition
)

func (t tokenType) String() string {
	return [...]string{"SingleCharacter", "EscapeSequence", "Repetition"}[t]
}

type token struct {
	kind  tokenType
	token string
}

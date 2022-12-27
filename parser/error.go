package parser

import "fmt"

const (
	Continue = 1
	Fatal    = 2
)

var (
	ParseContinue = ParseError{kind: Continue}
	ParseFatal    = ParseError{kind: Fatal}
)

type IParseError interface {
	error
	IsContinue() bool
	IsFatal() bool
}

type ParseError struct {
	kind   int
	reason string
}

func (p ParseError) Error() string {
	if p.kind == Continue {
		return fmt.Sprintf("parse fail need more data: %s", p.reason)
	}
	return fmt.Sprintf("parse fatal error: %s", p.reason)
}
func (p ParseError) Is(t error) bool {
	e, ok := t.(ParseError)
	if !ok {
		return false
	}
	return e.kind == p.kind
}
func (p ParseError) WithReason(reason string) ParseError {
	t := p
	t.reason = reason
	return t
}

func (p ParseError) IsFatal() bool {
	return p.kind == Fatal
}

func (p ParseError) IsContinue() bool {
	return p.kind == Continue
}

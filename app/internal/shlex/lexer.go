package shlex

import (
	"strings"
)

type Lexer struct {
	args          []string
	current       strings.Builder
	inSingleQuote bool
	inDoubleQuote bool
	escaped       bool
	inArg         bool
}

func Split(line string) ([]string, error) {
	l := &Lexer{}
	return l.Run(line)
}

func (l *Lexer) Run(line string) ([]string, error) {
	for _, r := range line {
		if l.escaped {
			l.handleEscaped(r)
			continue
		}

		switch {
		case l.inSingleQuote:
			l.handleSingleQuote(r)
		case l.inDoubleQuote:
			l.handleDoubleQuote(r)
		case r == '\\':
			l.escaped = true
			l.inArg = true
		case r == '\'':
			l.inSingleQuote = true
			l.inArg = true
		case r == '"':
			l.inDoubleQuote = true
			l.inArg = true
		case isSpace(r):
			l.flushArg()
		default:
			l.current.WriteRune(r)
			l.inArg = true
		}
	}

	return l.finish()
}

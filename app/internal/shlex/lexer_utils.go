package shlex

import (
	"fmt"
)

func (l *Lexer) handleEscaped(r rune) {
	if l.inDoubleQuote {
		switch r {
		case '"', '\\':
			l.current.WriteRune(r)
		default:
			l.current.WriteRune('\\')
			l.current.WriteRune(r)
		}
	} else {
		l.current.WriteRune(r)
	}
	l.inArg = true
	l.escaped = false
}

func (l *Lexer) handleSingleQuote(r rune) {
	switch r {
	case '\'':
		l.inSingleQuote = false
	default:
		l.current.WriteRune(r)
	}
	l.inArg = true
}

func (l *Lexer) handleDoubleQuote(r rune) {
	switch r {
	case '"':
		l.inDoubleQuote = false
	case '\\':
		l.escaped = true
	default:
		l.current.WriteRune(r)
	}
	l.inArg = true
}

func (l *Lexer) flushArg() {
	if l.inArg {
		l.args = append(l.args, l.current.String())
		l.current.Reset()
		l.inArg = false
	}
}

func (l *Lexer) finish() ([]string, error) {
	if l.inSingleQuote {
		return nil, fmt.Errorf("unmatched single quote")
	}
	if l.inDoubleQuote {
		return nil, fmt.Errorf("unmatched double quote")
	}
	if l.escaped {
		return nil, fmt.Errorf("trailing backslash")
	}

	l.flushArg()
	return l.args, nil
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

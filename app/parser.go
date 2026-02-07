package main

import (
	"fmt"
	"strings"
)

func parseArgs(line string) ([]string, error) {
	var args []string
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	inArg := false

	for _, r := range line {
		switch {
		case inSingleQuote:
			if r == '\'' {
				inSingleQuote = false
			} else {
				current.WriteRune(r)
			}
			inArg = true
		case inDoubleQuote:
			if r == '"' {
				inDoubleQuote = false
			} else {
				current.WriteRune(r)
			}
			inArg = true
		case r == '\'':
			inSingleQuote = true
			inArg = true
		case r == '"':
			inDoubleQuote = true
			inArg = true
		case r == ' ' || r == '\t' || r == '\n' || r == '\r':
			if inArg {
				args = append(args, current.String())
				current.Reset()
				inArg = false
			}
		default:
			current.WriteRune(r)
			inArg = true
		}
	}

	if inSingleQuote {
		return nil, fmt.Errorf("unmatched single quote")
	}
	if inDoubleQuote {
		return nil, fmt.Errorf("unmatched double quote")
	}

	if inArg {
		args = append(args, current.String())
	}
	return args, nil
}

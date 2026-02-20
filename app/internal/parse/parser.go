package parse

import "fmt"

type ParseResult struct {
	Cmd           string
	Args          []string
	Redirect      *Redirect
	ErrorRedirect *Redirect
	Err           error
}

type Pipeline struct {
	Commands []*ParseResult
}

func parseSingleCommand(tokens []string) *ParseResult {
	if len(tokens) == 0 {
		return &ParseResult{Err: fmt.Errorf("no command provided")}
	}

	result := &ParseResult{
		Cmd:  tokens[0],
		Args: []string{},
	}

	missingTargetErr := func(token string) {
		result.Err = fmt.Errorf("syntax error: redirect operator %q at end of line with no target", token)
	}

	for i := 1; i < len(tokens); i++ {
		token := tokens[i]
		switch token {
		case ">", "1>":
			if i+1 >= len(tokens) {
				missingTargetErr(token)
				return result
			}
			target := tokens[i+1]
			result.Redirect = &Redirect{
				FD:     1,
				Target: target,
				Append: false,
			}
			i++
		case ">>", "1>>":
			if i+1 >= len(tokens) {
				missingTargetErr(token)
				return result
			}
			target := tokens[i+1]
			result.Redirect = &Redirect{
				FD:     1,
				Target: target,
				Append: true,
			}
			i++
		case "2>":
			if i+1 >= len(tokens) {
				missingTargetErr(token)
				return result
			}
			target := tokens[i+1]
			result.ErrorRedirect = &Redirect{
				FD:     2,
				Target: target,
				Append: false,
			}
			i++
		case "2>>":
			if i+1 >= len(tokens) {
				missingTargetErr(token)
				return result
			}
			target := tokens[i+1]
			result.ErrorRedirect = &Redirect{
				FD:     2,
				Target: target,
				Append: true,
			}
			i++
		default:
			result.Args = append(result.Args, token)
		}
	}

	return result
}

func ParseCommand(tokens []string) *Pipeline {
	commands := SplitByPipe(tokens)
	pipeline := &Pipeline{
		Commands: make([]*ParseResult, len(commands)),
	}

	for i, cmdTokens := range commands {
		pipeline.Commands[i] = parseSingleCommand(cmdTokens)
	}

	return pipeline
}

func SplitByPipe(tokens []string) [][]string {
	var result [][]string
	var current []string

	for _, token := range tokens {
		if token == "|" {
			result = append(result, current)
			current = []string{}
		} else {
			current = append(current, token)
		}
	}
	result = append(result, current)

	return result
}

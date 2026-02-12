package parse

import "fmt"

type ParseResult struct {
	Cmd           string
	Args          []string
	Redirect      *Redirect
	ErrorRedirect *Redirect
	Err           error
}

func ParseCommand(tokens []string) *ParseResult {
	if len(tokens) == 0 {
		return &ParseResult{Err: fmt.Errorf("no command provided")}
	}

	result := &ParseResult{
		Cmd:  tokens[0],
		Args: []string{},
	}

	for i := 1; i < len(tokens); i++ {
		token := tokens[i]
		switch token {
		case ">", "1>":
			if i+1 >= len(tokens) {
				result.Err = fmt.Errorf("syntax error: redirect operator %q at end of line with no target", token)
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
				result.Err = fmt.Errorf("syntax error: redirect operator %q at end of line with no target", token)
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
				result.Err = fmt.Errorf("syntax error: redirect operator %q at end of line with no target", token)
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
				result.Err = fmt.Errorf("syntax error: redirect operator %q at end of line with no target", token)
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

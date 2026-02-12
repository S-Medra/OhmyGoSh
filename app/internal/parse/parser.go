package parse

import "fmt"

func ParseCommand(tokens []string) (string, []string, *Redirect, *Redirect, error) {
	if len(tokens) == 0 {
		return "", nil, nil, nil, fmt.Errorf("no command provided")
	}

	cmd := tokens[0]
	var args []string
	var redirect *Redirect
	var errorRedirect *Redirect

	for i := 1; i < len(tokens); i++ {
		token := tokens[i]
		switch token {
		case ">", "1>":
			if i+1 >= len(tokens) {
				return "", nil, nil, nil, fmt.Errorf("syntax error: redirect operator %q at end of line with no target", token)
			}
			target := tokens[i+1]
			redirect = &Redirect{
				FD:     1,
				Target: target,
				Append: false,
			}
			i++
		case ">>", "1>>":
			if i+1 >= len(tokens) {
				return "", nil, nil, nil, fmt.Errorf("syntax error: redirect operator %q at end of line with no target", token)
			}
			target := tokens[i+1]
			redirect = &Redirect{
				FD:     1,
				Target: target,
				Append: true,
			}
			i++
		case "2>":
			if i+1 >= len(tokens) {
				return "", nil, nil, nil, fmt.Errorf("syntax error: redirect operator %q at end of line with no target", token)
			}
			target := tokens[i+1]
			errorRedirect = &Redirect{
				FD:     2,
				Target: target,
				Append: false,
			}
			i++
		case "2>>":
			if i+1 >= len(tokens) {
				return "", nil, nil, nil, fmt.Errorf("syntax error: redirect operator %q at end of line with no target", token)
			}
			target := tokens[i+1]
			errorRedirect = &Redirect{
				FD:     2,
				Target: target,
				Append: true,
			}
			i++
		default:
			args = append(args, token)
		}
	}

	return cmd, args, redirect, errorRedirect, nil
}

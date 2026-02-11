package parse

import "fmt"

func ParseCommand(tokens []string) (string, []string, *Redirect, error) {
	if len(tokens) == 0 {
		return "", nil, nil, fmt.Errorf("no command provided")
	}

	cmd := tokens[0]
	var args []string
	var redirect *Redirect

	for i := 1; i < len(tokens); i++ {
		token := tokens[i]
		switch token {
		case ">", "1>":
			if i+1 >= len(tokens) {
				return "", nil, nil, fmt.Errorf("syntax error: redirect operator %q at end of line with no target", token)
			}
			target := tokens[i+1]
			fd := 1
			redirect = &Redirect{
				FD:     fd,
				Target: target,
			}
			i++
		default:
			args = append(args, token)
		}
	}

	return cmd, args, redirect, nil
}

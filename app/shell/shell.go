package shell

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ergochat/readline"
	"github.com/ixiSam/OhmyGoSh/app/internal/parse"
	"github.com/ixiSam/OhmyGoSh/app/internal/shlex"
)

type CommandFunc func(args []string, out io.Writer) error

type Shell struct {
	cwd      string
	out      io.Writer
	err      io.Writer
	rl       *readline.Instance
	commands map[string]CommandFunc
}

func New(out, errW io.Writer) (*Shell, error) {
	cwd, _ := os.Getwd()
	s := &Shell{
		cwd: cwd,
		out: out,
		err: errW,
	}

	s.commands = defaultBuiltins(s)

	// Pass a live reference to the commands map so the completer
	// stays in sync if commands are added/removed later.
	completer := newCompleter(func() map[string]CommandFunc {
		return s.commands
	})

	rl, err := readline.NewFromConfig(&readline.Config{
		Prompt:            "OhmyGoSh$ ",
		AutoComplete:      completer,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
	})
	if err != nil {
		return nil, fmt.Errorf("initializing readline: %w", err)
	}
	s.rl = rl

	return s, nil
}

// Close releases terminal resources held by readline.
func (s *Shell) Close() error {
	return s.rl.Close()
}

func (s *Shell) Run() error {
	for {
		line, err := s.rl.ReadLine()
		if err != nil {
			if err == readline.ErrInterrupt {
				continue // Ctrl+C: clear line, keep going
			}
			// io.EOF (Ctrl+D) or unexpected error
			return nil
		}

		line = strings.TrimSpace(line)
		tokens, err := shlex.Split(line)
		if err != nil {
			fmt.Fprintln(s.err, "Error:", err)
			continue
		}

		if len(tokens) == 0 {
			continue
		}

		result := parse.ParseCommand(tokens)
		if result.Err != nil {
			fmt.Fprintln(s.err, "Error:", result.Err)
			continue
		}

		var rp redirectPair
		if err := rp.open(result.Redirect, result.ErrorRedirect); err != nil {
			fmt.Fprintln(s.err, "Error:", err)
			continue
		}

		// Close redirect files at end of each iteration to avoid leaking file handles.
		// We use an immediately invoked function to ensure Close() runs even if we continue early.
		func() {
			defer rp.Close()

			cmdOut := rp.stdout.Writer(s.out)
			cmdErr := rp.stderr.Writer(s.err)

			if cmd, ok := s.commands[result.Cmd]; ok {
				if err := cmd(result.Args, cmdOut); err != nil {
					fmt.Fprintln(s.err, "Error:", err)
				}
				return
			}

			if err := runExternal(result.Cmd, result.Args, cmdOut, cmdErr); err != nil {
				fmt.Fprintln(s.err, "Error:", err)
			}
		}()
	}
}

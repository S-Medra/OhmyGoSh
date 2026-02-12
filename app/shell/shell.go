package shell

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/ixiSam/OhmyGoSh/app/internal/parse"
	"github.com/ixiSam/OhmyGoSh/app/internal/shlex"
)

type CommandFunc func(args []string, out io.Writer) error

type Shell struct {
	out      io.Writer
	err      io.Writer
	reader   *bufio.Reader
	commands map[string]CommandFunc
}

func New(in io.Reader, out, err io.Writer) *Shell {
	s := &Shell{
		out:    out,
		err:    err,
		reader: bufio.NewReader(in),
	}

	s.commands = defaultBuiltins(s)

	return s
}

func (s *Shell) Run() error {
	for {
		fmt.Fprint(s.out, "OhmyGoSh$ ")

		line, err := s.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			fmt.Fprintln(s.err, "Error reading command:", err)
			return err
		}

		line = strings.TrimSuffix(line, "\n")
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
		defer rp.Close()

		cmdOut := rp.stdout.Writer(s.out)
		cmdErr := rp.stderr.Writer(s.err)

		if cmdFn, ok := s.commands[result.Cmd]; ok {
			if err := cmdFn(result.Args, cmdOut); err != nil {
				fmt.Fprintln(s.err, "Error:", err)
			}
			continue
		}

		if err := runExternal(result.Cmd, result.Args, cmdOut, cmdErr); err != nil {
			fmt.Fprintln(s.err, "Error:", err)
		}
	}
}

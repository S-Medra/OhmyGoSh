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
	in       io.Reader
	out      io.Writer
	err      io.Writer
	reader   *bufio.Reader
	commands map[string]CommandFunc
}

func New(in io.Reader, out, err io.Writer) *Shell {
	s := &Shell{
		in:     in,
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

		cmd, cmdArgs, redirect, errorRedirect, err := parse.ParseCommand(tokens)
		if err != nil {
			fmt.Fprintln(s.err, "Error:", err)
			continue
		}

		var rp redirectPair
		if err := rp.open(redirect, errorRedirect); err != nil {
			fmt.Fprintln(s.err, "Error:", err)
			continue
		}
		defer rp.Close()

		cmdOut := rp.stdout.Writer(s.out)
		cmdErr := rp.stderr.Writer(s.err)

		if cmdFn, ok := s.commands[cmd]; ok {
			if err := cmdFn(cmdArgs, cmdOut); err != nil {
				fmt.Fprintln(s.err, "Error:", err)
			}
			continue
		}

		if err := runExternal(cmd, cmdArgs, cmdOut, cmdErr); err != nil {
			fmt.Fprintln(s.err, "Error:", err)
		}
	}
}

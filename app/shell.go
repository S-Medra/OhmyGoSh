package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/S-Medra/OhmyGoSh/app/internal/shlex"
)

type CommandFunc func(args []string) error

type Shell struct {
	in       io.Reader
	out      io.Writer
	err      io.Writer
	reader   *bufio.Reader
	commands map[string]CommandFunc
}

func NewShell(in io.Reader, out, err io.Writer) *Shell {
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
				return nil // Exit gracefully on Ctrl+D
			}
			fmt.Fprintln(s.err, "Error reading command:", err)
			return err
		}

		line = strings.TrimSuffix(line, "\n")
		args, err := shlex.Split(line)
		if err != nil {
			fmt.Fprintln(s.err, "Error:", err)
			continue
		}

		if len(args) == 0 {
			continue
		}

		name := args[0]
		cmdArgs := args[1:]

		if cmdFn, ok := s.commands[name]; ok {
			if err := cmdFn(cmdArgs); err != nil {
				fmt.Fprintln(s.err, "Error:", err)
			}
			continue
		}

		if err := runExternal(name, cmdArgs, s.in, s.out, s.err); err != nil {
			fmt.Fprintln(s.err, "Error:", err)
		}
	}
}

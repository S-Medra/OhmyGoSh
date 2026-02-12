package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/S-Medra/OhmyGoSh/app/internal/parse"
	"github.com/S-Medra/OhmyGoSh/app/internal/shlex"
)

type CommandFunc func(args []string, out io.Writer) error

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

		rw, err := redirect.Open()
		if err != nil {
			fmt.Fprintf(s.err, "Error creating redirect file: %v\n", err)
			continue
		}
		defer rw.Close()

		erw, err := errorRedirect.Open()
		if err != nil {
			fmt.Fprintf(s.err, "Error creating error redirect file: %v\n", err)
			continue
		}
		defer erw.Close()

		cmdOut := rw.Writer(s.out)
		cmdErr := erw.Writer(s.err)

		if cmdFn, ok := s.commands[cmd]; ok {
			if err := cmdFn(cmdArgs, cmdOut); err != nil {
				fmt.Fprintln(s.err, "Error:", err)
			}
			continue
		}

		if err := runExternal(cmd, cmdArgs, s.in, cmdOut, cmdErr); err != nil {
			fmt.Fprintln(s.err, "Error:", err)
		}
	}
}

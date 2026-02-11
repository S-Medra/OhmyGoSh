package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func defaultBuiltins(s *Shell) map[string]CommandFunc {
	return map[string]CommandFunc{
		"exit": s.exitCmd,
		"echo": s.echoCmd,
		"type": s.typeCmd,
		"pwd":  s.pwdCmd,
		"cd":   s.cdCmd,
	}
}

func (s *Shell) exitCmd(args []string, out io.Writer) error {
	if len(args) == 0 {
		os.Exit(0)
	}

	status, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintln(s.err, "exit: numeric argument required")
		os.Exit(1)
	}

	os.Exit(status)
	return nil
}

func (s *Shell) echoCmd(args []string, out io.Writer) error {
	output := strings.Join(args, " ")
	fmt.Fprintln(out, output)
	return nil
}

func (s *Shell) typeCmd(args []string, out io.Writer) error {
	if len(args) == 0 {
		fmt.Fprintln(s.err, "Received no args")
		return nil
	}

	name := args[0]

	if _, ok := s.commands[name]; ok {
		fmt.Fprintln(out, name+" is a shell builtin")
		return nil
	}

	if path, err := exec.LookPath(name); err == nil {
		fmt.Fprintln(out, name+" is "+path)
		return nil
	}

	fmt.Fprintln(s.err, name+" not found")
	return nil
}

func (s *Shell) pwdCmd(args []string, out io.Writer) error {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(s.err, "Error getting dir: %v\n", err)
		return err
	}
	fmt.Fprintln(out, pwd)
	return nil
}

func (s *Shell) cdCmd(args []string, out io.Writer) error {
	var target string
	if len(args) == 0 || args[0] == "~" {
		var err error
		target, err = os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(s.err, "cd: %v\n", err)
			return err
		}
	} else {
		target = args[0]
	}
	err := os.Chdir(target)
	if err != nil {
		fmt.Fprintf(s.err, "cd: %s: no such file or directory\n", target)
	}
	return nil
}

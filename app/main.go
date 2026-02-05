package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type commandFunc func(args []string) error

var commands map[string]commandFunc

func main() {
	commands = map[string]commandFunc{
		"exit": exitCmd,
		"echo": echoCmd,
		"type": typeCmd,
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading command:", err)
			return
		}

		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		fields := strings.Fields(trimmed)
		if len(fields) == 0 {
			continue
		}

		name := fields[0]
		args := fields[1:]

		if cmdFn, ok := commands[name]; ok {
			if err := cmdFn(args); err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
			}
			continue
		}

		fmt.Println(name + ": not found")
	}
}

func exitCmd(args []string) error {
	if len(args) == 0 {
		os.Exit(0)
	}

	status, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "exit: numeric argument required")
		os.Exit(1)
	}

	os.Exit(status)
	return nil
}

func echoCmd(args []string) error {
	fmt.Println(strings.Join(args, " "))
	return nil
}

func typeCmd(args []string) error {
	if len(args) == 0 {
		fmt.Println("Received no args")
		return nil
	}

	name := args[0]

	if _, ok := commands[name]; ok {
		fmt.Println(name + " is a shell builtin")
		return nil
	}

	if path, err := exec.LookPath(name); err == nil {
		fmt.Println(name + " is " + path)
		return nil
	}

	fmt.Println(name + " not found")
	return nil
}

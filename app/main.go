package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		fmt.Print("$ ")
		cmd, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading command:", err)
			return
		}
		if cmd == "exit\n" {
			return
		}
		if strings.HasPrefix(cmd, "echo ") {
			input := cmd[len("echo ") : len(cmd)-1]
			fmt.Println(input)
			continue
		}

		if strings.HasPrefix(cmd, "type ") {
			arg := strings.TrimSpace(cmd[len("type "):])

			switch arg {
			case "exit", "echo", "type":
				fmt.Println(arg + " is a shell builtin")
			default:
				fmt.Println(arg + " not found")
			}
			continue
		}
		fmt.Println(cmd[:len(cmd)-1] + ": not found")
	}
}

package main

import (
	"fmt"
	"os"

	"github.com/ixiSam/OhmyGoSh/app/shell"
)

func main() {
	s := shell.New(os.Stdin, os.Stdout, os.Stderr)

	if err := s.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Shell exited with error:", err)
		os.Exit(1)
	}
}

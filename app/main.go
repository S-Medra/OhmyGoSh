package main

import (
	"fmt"
	"os"

	"github.com/ixiSam/OhmyGoSh/app/shell"
)

func main() {
	s, err := shell.New(os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to initialize shell:", err)
		os.Exit(1)
	}
	defer s.Close()

	if err := s.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Shell exited with error:", err)
		os.Exit(1)
	}
}

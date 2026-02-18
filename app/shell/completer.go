package shell

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// implements the readline.AutoCompleter interface.
type shellCompleter struct {
	commands func() map[string]CommandFunc

	mu        sync.Mutex
	pathCache []string
}

func newCompleter(commandsFn func() map[string]CommandFunc) *shellCompleter {
	c := &shellCompleter{commands: commandsFn}
	go c.scanPath()
	return c
}

// Do is required by the readline.AutoCompleter interface.
func (c *shellCompleter) Do(line []rune, pos int) ([][]rune, int) {
	lineStr := string(line[:pos])

	words := strings.Fields(lineStr)
	trailingSpace := len(lineStr) > 0 && lineStr[len(lineStr)-1] == ' '

	// Command position: complete command names.
	if len(words) == 0 || (len(words) == 1 && !trailingSpace) {
		prefix := ""
		if len(words) == 1 {
			prefix = words[0]
		}
		return c.completeCommand(prefix)
	}

	// Argument position: complete file/directory paths.
	prefix := ""
	if !trailingSpace {
		prefix = words[len(words)-1]
	}

	return c.completePath(prefix)
}

// completeCommand returns candidates matching the given prefix from
// builtins and PATH executables.
func (c *shellCompleter) completeCommand(prefix string) ([][]rune, int) {
	var candidates []string

	for name := range c.commands() {
		if strings.HasPrefix(name, prefix) {
			candidates = append(candidates, name)
		}
	}

	c.mu.Lock()
	snapshot := c.pathCache
	c.mu.Unlock()

	for _, name := range snapshot {
		if strings.HasPrefix(name, prefix) {
			candidates = append(candidates, name)
		}
	}

	return toRuneCandidates(candidates, prefix)
}

// completePath returns file/directory name candidates matching a partial path.
func (c *shellCompleter) completePath(prefix string) ([][]rune, int) {
	dir := "."
	base := prefix

	if i := strings.LastIndex(prefix, "/"); i >= 0 {
		dir = prefix[:i+1]
		base = prefix[i+1:]
	}

	if dir == "" {
		dir = "/"
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, 0
	}

	var candidates []string
	for _, entry := range entries {
		name := entry.Name()
		// Hide dotfiles unless the user explicitly typed a dot.
		if strings.HasPrefix(name, ".") && !strings.HasPrefix(base, ".") {
			continue
		}
		if !strings.HasPrefix(name, base) {
			continue
		}
		candidate := name
		if dir != "." {
			candidate = dir + name
		}
		if entry.IsDir() {
			candidate += "/"
		}
		candidates = append(candidates, candidate)
	}

	return toRuneCandidates(candidates, prefix)
}

// scanPath scans all directories in $PATH and populates the cache
// with executable names. Designed to be called in a goroutine.
func (c *shellCompleter) scanPath() {
	seen := make(map[string]struct{})
	var result []string

	pathEnv := os.Getenv("PATH")
	for _, dir := range filepath.SplitList(pathEnv) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if _, ok := seen[name]; ok {
				continue
			}
			info, err := entry.Info()
			if err != nil {
				continue
			}
			if info.Mode()&0111 != 0 {
				seen[name] = struct{}{}
				result = append(result, name)
			}
		}
	}

	c.mu.Lock()
	c.pathCache = result
	c.mu.Unlock()
}

// toRuneCandidates converts string candidates into the [][]rune, int format
// required by the readline.AutoCompleter interface.
func toRuneCandidates(candidates []string, prefix string) ([][]rune, int) {
	if len(candidates) == 0 {
		return nil, 0
	}

	result := make([][]rune, 0, len(candidates))
	for _, c := range candidates {
		suffix := c[len(prefix):]
		result = append(result, []rune(suffix+" "))
	}
	return result, len(prefix)
}

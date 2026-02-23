# OhmyGoSh

A custom Unix shell implementation in Go, built as a learning project to understand how shells work under the hood.

## Features

- **Built-in commands**: `exit`, `echo`, `type`, `pwd`, `cd` (with `~` expansion)
- **Pipes**: Chain commands with `|`
- **Redirection**: stdout (`>`, `>>`) and stderr (`2>`, `2>>`)
- **External commands**: Execute any binary in PATH
- **Tab completion**: Builtins, PATH executables, and file paths
- **Interactive**: Ctrl+C interrupt handling, command history with arrow keys

## Architecture

```
app/
├── main.go              # Entry point
├── shell/
│   ├── shell.go         # Core loop
│   ├── builtins.go      # Builtin commands
│   ├── pipeline.go      # Pipe execution
│   ├── executor.go      # Command excution
│   ├── external.go      # External commands
│   └── completer.go     # Tab completion
└── internal/
    ├── parse/           # Parsing & redirection
    └── shlex/           # Lexer & tokenization
```

Uses Go's `internal` package pattern to enforce encapsulation.

## Design Decisions

- **Modularity**: Single-responsibility files - each handles one concern.
- **Extendability**: For example Builtins in `map[string]CommandFunc` - add new commands in one line.
- **Testability**: Used tests for parse and shlex packages to ensure correctness and prevent regressions.

## Development Journey

Started as a simple REPL, added features one at a time.

- **Basic REPL** - print prompt, read line, print output, repeat. 
*Learned how a read-print loop works, and the difference between handling Ctrl+C (interrupt) vs Ctrl+D (exit).*

- **Builtins** - implemented exit, echo, type. 
*Learned that commands run in-process, unlike external binaries.*

- **Dispatch map** - refactored builtins into `map[string]CommandFunc`. 
*This gives O(1) lookup and makes adding new commands trivial.*

- **External commands** - added os/exec to run PATH binaries. 
*Learned how exec.LookPath searches $PATH, and that external commands fork automatically via exec.Command.*

- **Lexer** - built from scratch to handle quotes and escaping. 
*Learned that tokenizing requires a state machine to track quote contexts.*

- **Redirection** - stdout with >, >> and stderr with 2>. 
*Used file descriptors 1 and 2 for stdout/stderr, learned how O_CREATE/O_APPEND flags control file behavior.*

- **Tab completion** - completers for builtins, PATH, and file paths. 
*Built AutoCompleter to be context aware (command vs path), made sure to do background scanning using sync.Mutex for thread safety.*

- **Pipes** - chained commands with |. 
*Learned that os.Pipe() creates connected file descriptors, and closing write end signals EOF to reader.*

- **Piped builtins** - pipelines with builtin commands. 
*Cause builtins run in the same process they can't block on stdin - goroutines solve this by running them concurrently.*

## Quick Start

```bash
go run ./app
```
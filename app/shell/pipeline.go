package shell

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/ixiSam/OhmyGoSh/app/internal/parse"
)

type Pipeline struct {
	commands []*parse.ParseResult
	stdout   io.Writer
	stderr   io.Writer
	builtins map[string]CommandFunc

	runningCmds []*exec.Cmd
	builtinErrs []error
	pipes       [][2]*os.File
}

func NewPipeline(cmds []*parse.ParseResult, stdout, stderr io.Writer, builtins map[string]CommandFunc) *Pipeline {
	return &Pipeline{
		commands: cmds,
		stdout:   stdout,
		stderr:   stderr,
		builtins: builtins,
	}
}

func (p *Pipeline) Run() error {
	if len(p.commands) == 0 {
		return nil
	}

	if len(p.commands) == 1 {
		return p.runSingle()
	}

	if err := p.createPipes(); err != nil {
		return err
	}

	return p.runPiped()
}

func (p *Pipeline) runSingle() error {
	cmd := p.commands[0]
	if builtin, ok := p.builtins[cmd.Cmd]; ok {
		return builtin(cmd.Args, p.stdout)
	}
	return runExternal(cmd.Cmd, cmd.Args, p.stdout, p.stderr)
}

func (p *Pipeline) createPipes() error {
	count := len(p.commands) - 1
	p.pipes = make([][2]*os.File, count)
	for i := range count {
		readEnd, writeEnd, err := os.Pipe()
		if err != nil {
			return fmt.Errorf("creating pipe %d: %v", i, err)
		}
		p.pipes[i] = [2]*os.File{readEnd, writeEnd}
	}
	return nil
}

func (p *Pipeline) runPiped() error {
	p.runningCmds = make([]*exec.Cmd, len(p.commands))
	p.builtinErrs = make([]error, len(p.commands))

	for i, cmd := range p.commands {
		stdin, stdoutW := p.getIO(i)

		if builtin, ok := p.builtins[cmd.Cmd]; ok {
			go func(idx int, b CommandFunc, args []string, out io.Writer) {
				p.builtinErrs[idx] = b(args, out)
			}(i, builtin, cmd.Args, stdoutW)
			continue
		}

		execCmd := exec.Command(cmd.Cmd, cmd.Args...)
		execCmd.Stdin = stdin
		execCmd.Stdout = stdoutW
		execCmd.Stderr = p.stderr

		if err := execCmd.Start(); err != nil {
			return fmt.Errorf("%s: failed to start: %v", cmd.Cmd, err)
		}
		p.runningCmds[i] = execCmd
	}

	p.cleanupPipes()
	p.waitForExternals()
	return p.checkErrors()
}

func (p *Pipeline) getIO(index int) (io.Reader, io.Writer) {
	var stdin io.Reader
	var stdoutW io.Writer

	if index == 0 {
		stdin = os.Stdin
	} else {
		stdin = p.pipes[index-1][0]
	}

	if index == len(p.commands)-1 {
		stdoutW = p.stdout
	} else {
		stdoutW = p.pipes[index][1]
	}

	return stdin, stdoutW
}

func (p *Pipeline) cleanupPipes() {
	for _, pipe := range p.pipes {
		pipe[1].Close()
		pipe[0].Close()
	}
}

func (p *Pipeline) waitForExternals() {
	for _, runningCmd := range p.runningCmds {
		if runningCmd != nil {
			runningCmd.Wait()
		}
	}
}

func (p *Pipeline) checkErrors() error {
	for _, err := range p.builtinErrs {
		if err != nil {
			return err
		}
	}
	return nil
}

package shell

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/ixiSam/OhmyGoSh/app/internal/parse"
)

func runPipeline(pipeline *parse.Pipeline, stdout, stderr io.Writer) error {
	cmds := pipeline.Commands
	if len(cmds) == 0 {
		return nil
	}

	if len(cmds) == 1 {
		return runExternal(cmds[0].Cmd, cmds[0].Args, stdout, stderr)
	}

	pipes, err := createPipes(len(cmds) - 1)
	if err != nil {
		return err
	}

	return runPipedCommands(cmds, pipes, stdout, stderr)
}

func createPipes(count int) ([][2]*os.File, error) {
	pipes := make([][2]*os.File, count)
	for i := range count {
		readEnd, writeEnd, err := os.Pipe()
		if err != nil {
			return nil, fmt.Errorf("creating pipe %d: %v", i, err)
		}
		pipes[i] = [2]*os.File{readEnd, writeEnd}
	}
	return pipes, nil
}

func runPipedCommands(commands []*parse.ParseResult, pipes [][2]*os.File, stdout, stderr io.Writer) error {
	runningCmds := make([]*exec.Cmd, len(commands))

	for i, cmd := range commands {
		var stdin io.Reader
		var stdoutW io.Writer

		if i == 0 {
			stdin = os.Stdin
		} else {
			stdin = pipes[i-1][0]
		}

		if i == len(commands)-1 {
			stdoutW = stdout
		} else {
			stdoutW = pipes[i][1]
		}

		execCmd := exec.Command(cmd.Cmd, cmd.Args...)
		execCmd.Stdin = stdin
		execCmd.Stdout = stdoutW
		execCmd.Stderr = stderr

		if err := execCmd.Start(); err != nil {
			return fmt.Errorf("%s: failed to start: %v", cmd.Cmd, err)
		}
		runningCmds[i] = execCmd
	}

	for _, pipe := range pipes {
		pipe[1].Close()
		pipe[0].Close()
	}

	for _, runningCmd := range runningCmds {
		runningCmd.Wait()
	}

	return nil
}

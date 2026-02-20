package shell

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/ixiSam/OhmyGoSh/app/internal/parse"
)

type redirectPair struct {
	stdout *parse.RedirectWriter
	stderr *parse.RedirectWriter
}

func (rp *redirectPair) open(redirect, errorRedirect *parse.Redirect) error {
	rw, err := redirect.Open()
	if err != nil {
		return fmt.Errorf("creating redirect file: %v", err)
	}
	rp.stdout = rw

	erw, err := errorRedirect.Open()
	if err != nil {
		rw.Close()
		return fmt.Errorf("creating redirect file: %v", err)
	}
	rp.stderr = erw

	return nil
}

func (rp *redirectPair) Close() error {
	if rp.stdout != nil {
		if err := rp.stdout.Close(); err != nil {
			return err
		}
	}
	if rp.stderr != nil {
		if err := rp.stderr.Close(); err != nil {
			return err
		}
	}
	return nil
}

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
	for i := 0; i < count; i++ {
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

func runExternal(name string, args []string, out, err io.Writer) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = out
	cmd.Stderr = err

	if errRun := cmd.Run(); errRun != nil {
		var ee *exec.Error
		if errors.As(errRun, &ee) && ee.Err == exec.ErrNotFound {
			fmt.Fprintln(err, name+": not found")
			return nil
		}

		if _, ok := errRun.(*exec.ExitError); ok {
			return nil
		}

		return fmt.Errorf("%s: execution failed: %v", name, errRun)
	}

	return nil
}

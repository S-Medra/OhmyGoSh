package shell

import (
	"errors"
	"fmt"
	"io"
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
		return fmt.Errorf("creating error redirect file: %v", err)
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

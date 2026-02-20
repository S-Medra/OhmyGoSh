package shell

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
)

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

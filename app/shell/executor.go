package shell

import (
	"fmt"

	"github.com/ixiSam/OhmyGoSh/app/internal/parse"
)

type Redirects struct {
	stdout *parse.RedirectWriter
	stderr *parse.RedirectWriter
}

func (rp *Redirects) open(redirect, errorRedirect *parse.Redirect) error {
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

func (rp *Redirects) Close() error {
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

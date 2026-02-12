package parse

import (
	"io"
	"os"
)

type Redirect struct {
	FD     int
	Target string
	Append bool
}

type RedirectWriter struct {
	file *os.File
}

func (r *Redirect) Open() (*RedirectWriter, error) {
	if r == nil {
		return &RedirectWriter{file: nil}, nil
	}

	flag := os.O_WRONLY | os.O_CREATE
	if r.Append {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}

	file, err := os.OpenFile(r.Target, flag, 0644)
	if err != nil {
		return nil, err
	}

	return &RedirectWriter{file: file}, nil
}

func (rw *RedirectWriter) Writer(fallback io.Writer) io.Writer {
	if rw.file != nil {
		return rw.file
	}
	return fallback
}

func (rw *RedirectWriter) Close() error {
	if rw.file != nil {
		return rw.file.Close()
	}
	return nil
}

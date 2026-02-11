package parse

import (
	"io"
	"os"
)

type Redirect struct {
	FD     int
	Target string
}

type RedirectWriter struct {
	file *os.File
}

func (r *Redirect) Open() (*RedirectWriter, error) {
	if r == nil {
		return &RedirectWriter{file: nil}, nil
	}

	file, err := os.Create(r.Target)
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

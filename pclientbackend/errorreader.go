package pclientbackend

import (
	"io"
)

type ErrorReader struct{}

func (r *ErrorReader) Read(p []byte) (int, error) {
	return 0, io.EOF
}

func (r *ErrorReader) Write(p []byte) (int, error) {
	return len(p), nil
}

func (r *ErrorReader) Close() error {
	return nil
}

var errReader = new(ErrorReader)

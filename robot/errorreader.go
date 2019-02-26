package main

import (
	"io"
)

type ErrorReadWriter struct{}

func (r *ErrorReadWriter) Read(p []byte) (int, error) {
	return 0, io.EOF
}

func (r *ErrorReadWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (r *ErrorReadWriter) Close() error {
	return nil
}

var errReader = new(ErrorReadWriter)

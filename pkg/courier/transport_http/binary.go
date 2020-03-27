package transport_http

import (
	"bytes"
	"io"
)

func NewBinary() *Binary {
	return &Binary{}
}

var _ interface {
	io.Reader
	io.Writer
} = (*Binary)(nil)

// swagger:strfmt binary
type Binary struct {
	buf bytes.Buffer
}

func (binary *Binary) Read(p []byte) (n int, err error) {
	return binary.buf.Read(p)
}

func (binary *Binary) Write(p []byte) (n int, err error) {
	return binary.buf.Write(p)
}

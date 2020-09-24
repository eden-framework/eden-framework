package transport_http

import (
	"bytes"
	"io"

	"github.com/eden-framework/eden-framework/pkg/courier"
)

func NewFile(filename string, contentType string) *File {
	return &File{
		filename:    filename,
		contentType: contentType,
	}
}

var _ interface {
	io.Reader
	io.Writer
	IContentType
	courier.IMeta
} = (*File)(nil)

// swagger:strfmt binary
type File struct {
	filename    string
	contentType string
	buf         bytes.Buffer
}

func (file *File) Read(p []byte) (n int, err error) {
	return file.buf.Read(p)
}

func (file *File) Write(p []byte) (n int, err error) {
	return file.buf.Write(p)
}

func (file *File) Meta() courier.Metadata {
	metadata := courier.Metadata{}
	metadata.Add("Content-Disposition", "attachment; filename="+file.filename)
	return metadata
}

func (file *File) ContentType() string {
	return file.contentType
}

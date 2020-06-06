package httpx

import (
	"bytes"
)

func NewHTML(data []byte) *HTML {
	html := &HTML{}
	html.Write(data)
	return html
}

// swagger:strfmt html
type HTML struct {
	bytes.Buffer
}

func (h HTML) ContentType() string {
	return MIME_HTML
}

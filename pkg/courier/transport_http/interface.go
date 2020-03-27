package transport_http

import (
	"net/http"
)

type IMethod interface {
	Method() string
}

type IBytesGetter interface {
	Bytes() []byte
}

type IPath interface {
	Path() string
}

type IContentType interface {
	ContentType() string
}

type IStatus interface {
	Status() int
}

type IHttpRequestTransformer interface {
	TransformHttpRequest(req *http.Request)
}

package httpx

import (
	"net/http"
)

type MethodGet struct {
}

func (MethodGet) Method() string {
	return http.MethodGet
}

type MethodHead struct{}

func (MethodHead) Method() string {
	return http.MethodHead
}

type MethodPost struct{}

func (MethodPost) Method() string {
	return http.MethodPost
}

type MethodPut struct{}

func (MethodPut) Method() string {
	return http.MethodPut
}

type MethodPatch struct{}

func (MethodPatch) Method() string {
	return http.MethodPatch
}

type MethodDelete struct{}

func (MethodDelete) Method() string {
	return http.MethodDelete
}

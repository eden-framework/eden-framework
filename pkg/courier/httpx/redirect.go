package httpx

import (
	"fmt"
	"net/http"
	"strings"
)

type IRedirect interface {
	Redirect(host string) string
	Status() int
}

func RedirectWithStatusMovedPermanently(uri string) *statusMovedPermanently {
	return &statusMovedPermanently{
		location: location(uri),
	}
}

type statusMovedPermanently struct {
	location
}

func (s *statusMovedPermanently) Status() int {
	return http.StatusMovedPermanently
}

func RedirectWithStatusFound(uri string) *statusFound {
	return &statusFound{
		location: location(uri),
	}
}

type statusFound struct {
	location
}

func (s *statusFound) Status() int {
	return http.StatusFound
}

func RedirectWithStatusSeeOther(uri string) *statusSeeOther {
	return &statusSeeOther{
		location: location(uri),
	}
}

type statusSeeOther struct {
	location
}

func (s *statusSeeOther) Status() int {
	return http.StatusSeeOther
}

func RedirectWithStatusTemporaryRedirect(uri string) *statusTemporaryRedirect {
	return &statusTemporaryRedirect{
		location: location(uri),
	}
}

type statusTemporaryRedirect struct {
	location
}

func (s *statusTemporaryRedirect) Status() int {
	return http.StatusTemporaryRedirect
}

func RedirectWithStatusPermanentRedirect(uri string) *statusPermanentRedirect {
	return &statusPermanentRedirect{
		location: location(uri),
	}
}

type statusPermanentRedirect struct {
	location
}

func (s *statusPermanentRedirect) Status() int {
	return http.StatusPermanentRedirect
}

type location string

func (r *location) Error() string {
	return fmt.Sprintf("location : %s", r)
}

func (r *location) Redirect(host string) string {
	u := string(*r)
	if strings.HasPrefix(u, "/") {
		u = host + u
		if !strings.HasPrefix(u, "http") {
			u = "http://" + u
		}
	}
	return u
}

package httpx

import (
	"net/http"
)

type ICookie interface {
	Cookies() *http.Cookie
}

type WithCookie struct {
	c *http.Cookie
}

func (c *WithCookie) SetCookie(cookie *http.Cookie) {
	c.c = cookie
}

func (c *WithCookie) Cookies() *http.Cookie {
	return c.c
}

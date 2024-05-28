package HttpI

import (
	"io"
	"net/http"
)

type HttpI interface {
	Request(method Method, url string, headers Headers, cookies Cookies, body io.Reader) (*http.Response, error)
	SetProxy(rawUrl string) error
	SetCookies(rawUrl string, cookies Cookies)
	GetCookies(rawUrl string) Cookies
}

type Method string

const (
	GET     Method = "GET"
	POST    Method = "POST"
	PUT     Method = "PUT"
	HEAD    Method = "HEAD"
	DELETE  Method = "DELETE"
	OPTIONS Method = "OPTIONS"
)

type Headers map[string]string

func (H Headers) Set(key, value string) {
	H[key] = value
}

func (H Headers) Get(key string) string {
	return H[key]
}

type Cookies []*http.Cookie

func (C Cookies) Append(cookie *http.Cookie) Cookies {
	return append(C, cookie)
}

func (C Cookies) Get(name string) *http.Cookie {
	for _, cookie := range C {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

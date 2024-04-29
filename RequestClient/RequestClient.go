package RequestClient

import (
	fhttp "github.com/bogdanfinn/fhttp"
)

type RequestClient interface {
	Do(req *fhttp.Request) (*fhttp.Response, error)
	SetProxy(link string) error
}

package requestclient

import (
	fhttp "github.com/bogdanfinn/fhttp"
	"github.com/bogdanfinn/tls-client/profiles"
	"io"
	"sync"
)

type RequestClient interface {
	NewRequest(method, url string, body io.Reader) (*fhttp.Request, error)
	Do(req *fhttp.Request) (*fhttp.Response, error)
	SetProxy(link string) error
}

func init() {
	GetInstance()
}

var (
	Instance   *TlsClient
	clientOnce sync.Once
)

func GetInstance() *TlsClient {
	clientOnce.Do(func() {
		Instance = NewTlsClient(300, profiles.Okhttp4Android13)
	})
	return Instance
}

package RequestClient

import (
	"fmt"
	"github.com/aurorax-neo/go-logger"
	fhttp "github.com/bogdanfinn/fhttp"
	tlsClient "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"io"
)

type TlsClient struct {
	client tlsClient.HttpClient
}

func NewTlsClient(timeoutSeconds int, clientProfile profiles.ClientProfile) *TlsClient {
	jar := tlsClient.NewCookieJar()
	options := []tlsClient.HttpClientOption{
		tlsClient.WithTimeoutSeconds(timeoutSeconds),
		tlsClient.WithClientProfile(clientProfile),
		tlsClient.WithNotFollowRedirects(),
		tlsClient.WithCookieJar(jar),
	}
	client, err := tlsClient.NewHttpClient(tlsClient.NewNoopLogger(), options...)
	if err != nil {
		return nil
	}
	return &TlsClient{
		client: client,
	}
}

func (T *TlsClient) NewRequest(method, url string, body io.Reader) (*fhttp.Request, error) {
	request, err := fhttp.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	return request, nil

}

func (T *TlsClient) Do(req *fhttp.Request) (*fhttp.Response, error) {
	response, err := T.client.Do(req)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (T *TlsClient) SetProxy(link string) error {
	if link == "" {
		return nil
	}
	logger.Logger.Debug(fmt.Sprint("SetProxy: ", link))
	err := T.client.SetProxy(link)
	if err != nil {
		return err
	}
	return nil
}

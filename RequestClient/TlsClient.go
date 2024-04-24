package RequestClient

import (
	fhttp "github.com/bogdanfinn/fhttp"
	tlsClient "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"io"
	"math/rand"
	"time"
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

func RandomClientProfile() profiles.ClientProfile {
	// 初始化随机数生成器
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	clientProfiles := []profiles.ClientProfile{
		profiles.Firefox_102,
		profiles.Safari_15_6_1,
		profiles.Safari_16_0,
		profiles.Chrome_110,
		profiles.Okhttp4Android13,
		profiles.CloudflareCustom,
		profiles.Firefox_117,
	}
	// 随机选择一个
	randomIndex := rng.Intn(len(clientProfiles))
	return clientProfiles[randomIndex]
}

func NewRequest(method, url string, body io.Reader) (*fhttp.Request, error) {
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
	err := T.client.SetProxy(link)
	if err != nil {
		return err
	}
	return nil
}

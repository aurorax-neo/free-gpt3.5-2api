package TlsClient

import (
	"free-gpt3.5-2api/HttpI"
	"github.com/bogdanfinn/tls-client/profiles"
	"io"
	"net/http"
	"net/url"

	fhttp "github.com/bogdanfinn/fhttp"
	tlsClient "github.com/bogdanfinn/tls-client"
)

type TlsClient struct {
	Client    tlsClient.HttpClient
	ReqBefore handler
}

type handler func(req *fhttp.Request) error

func NewClient(timeoutSeconds int, clientProfile profiles.ClientProfile) *TlsClient {
	client, _ := tlsClient.NewHttpClient(tlsClient.NewNoopLogger(), []tlsClient.HttpClientOption{
		tlsClient.WithCookieJar(tlsClient.NewCookieJar()),
		tlsClient.WithNotFollowRedirects(),
		tlsClient.WithTimeoutSeconds(timeoutSeconds),
		tlsClient.WithClientProfile(clientProfile),
	}...)
	return &TlsClient{Client: client}
}

func convertResponse(resp *fhttp.Response) *http.Response {
	response := &http.Response{
		Status:           resp.Status,
		StatusCode:       resp.StatusCode,
		Proto:            resp.Proto,
		ProtoMajor:       resp.ProtoMajor,
		ProtoMinor:       resp.ProtoMinor,
		Header:           http.Header(resp.Header),
		Body:             resp.Body,
		ContentLength:    resp.ContentLength,
		TransferEncoding: resp.TransferEncoding,
		Close:            resp.Close,
		Uncompressed:     resp.Uncompressed,
		Trailer:          http.Header(resp.Trailer),
	}
	return response
}

func (TC *TlsClient) handleHeaders(req *fhttp.Request, headers HttpI.Headers) {
	if headers == nil {
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}

func (TC *TlsClient) handleCookies(req *fhttp.Request, cookies HttpI.Cookies) {
	if cookies == nil {
		return
	}
	for _, c := range cookies {
		req.AddCookie(&fhttp.Cookie{
			Name:       c.Name,
			Value:      c.Value,
			Path:       c.Path,
			Domain:     c.Domain,
			Expires:    c.Expires,
			RawExpires: c.RawExpires,
			MaxAge:     c.MaxAge,
			Secure:     c.Secure,
			HttpOnly:   c.HttpOnly,
			SameSite:   fhttp.SameSite(c.SameSite),
			Raw:        c.Raw,
			Unparsed:   c.Unparsed,
		})
	}
}

func (TC *TlsClient) Request(method HttpI.Method, url string, headers HttpI.Headers, cookies HttpI.Cookies, body io.Reader) (*http.Response, error) {
	req, err := fhttp.NewRequest(string(method), url, body)
	if err != nil {
		return nil, err
	}
	TC.handleHeaders(req, headers)
	TC.handleCookies(req, cookies)
	if TC.ReqBefore != nil {
		if err := TC.ReqBefore(req); err != nil {
			return nil, err
		}
	}
	do, err := TC.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return convertResponse(do), nil
}

func (TC *TlsClient) SetProxy(rawUrl string) error {
	return TC.Client.SetProxy(rawUrl)
}

func (TC *TlsClient) SetCookies(rawUrl string, cookies HttpI.Cookies) {
	if cookies == nil {
		return
	}
	u, err := url.Parse(rawUrl)
	if err != nil {
		return
	}
	var fCookies []*fhttp.Cookie
	for _, c := range cookies {
		fCookies = append(fCookies, &fhttp.Cookie{
			Name:       c.Name,
			Value:      c.Value,
			Path:       c.Path,
			Domain:     c.Domain,
			Expires:    c.Expires,
			RawExpires: c.RawExpires,
			MaxAge:     c.MaxAge,
			Secure:     c.Secure,
			HttpOnly:   c.HttpOnly,
			SameSite:   fhttp.SameSite(c.SameSite),
			Raw:        c.Raw,
			Unparsed:   c.Unparsed,
		})
	}
	TC.Client.GetCookieJar().SetCookies(u, fCookies)
}

func (TC *TlsClient) GetCookies(rawUrl string) HttpI.Cookies {
	currUrl, err := url.Parse(rawUrl)
	if err != nil {
		return nil
	}

	var cookies HttpI.Cookies
	for _, c := range TC.Client.GetCookies(currUrl) {
		cookies = append(cookies, &http.Cookie{
			Name:       c.Name,
			Value:      c.Value,
			Path:       c.Path,
			Domain:     c.Domain,
			Expires:    c.Expires,
			RawExpires: c.RawExpires,
			MaxAge:     c.MaxAge,
			Secure:     c.Secure,
			HttpOnly:   c.HttpOnly,
			SameSite:   http.SameSite(c.SameSite),
		})
	}
	return cookies
}

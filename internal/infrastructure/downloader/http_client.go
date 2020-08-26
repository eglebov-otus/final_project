package downloader

import (
	"image-previewer/internal/domain"
	"net/http"
	"net/url"
)

type Client interface {
	Get(rawURL string, headers domain.RequestHeaders) (resp *http.Response, err error)
}

type HTTPClient struct {
	client *http.Client
}

func (c *HTTPClient) Get(rawURL string, headers domain.RequestHeaders) (resp *http.Response, err error) {
	uri, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	if uri.Scheme == "" {
		uri.Scheme = "http"
	}

	//nolint:noctx
	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header(headers)

	return c.client.Do(req)
}

func NewHTTPClient(client *http.Client) *HTTPClient {
	return &HTTPClient{
		client: client,
	}
}

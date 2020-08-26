package downloader

import (
	"bytes"
	"image-previewer/internal/domain"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestHTTPClient_Get(t *testing.T) {
	t.Run("remote server responded with success status", func(t *testing.T) {
		client := NewTestClient(func(req *http.Request) *http.Response {
			require.Equal(t, "application/json", req.Header.Get("Accept"))
			require.Equal(t, "http://yandex.ru/image.png", req.URL.String())

			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		})

		httpClient := NewHTTPClient(client)

		headers := make(domain.RequestHeaders)
		headers["Accept"] = []string{"application/json"}

		resp, err := httpClient.Get("yandex.ru/image.png", headers)

		require.Nil(t, err)
		require.Equal(t, 200, resp.StatusCode)
	})
}

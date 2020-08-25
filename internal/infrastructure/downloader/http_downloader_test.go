package downloader

import (
	"bytes"
	"image-previewer/internal/domain/dto"
	"image-previewer/tests/mocks"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

//go:generate mockgen -destination=../../../tests/mocks/mock_http_client.go -package=mocks image-previewer/internal/infrastructure/downloader Client
//nolint:funlen
func TestHttpDownloader_Download(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Run("invalid remote server response status", func(t *testing.T) {
		client := mocks.NewMockClient(ctrl)
		client.
			EXPECT().
			Get(gomock.Any()).
			Return(&http.Response{
				StatusCode: http.StatusNotFound,
				Body:       ioutil.NopCloser(bytes.NewReader(nil)),
			}, nil)

		img, err := NewHTTPDownloader(client).Download(
			"http://yandex.ru/test.jpg",
			dto.ImageDimensions{
				Width:  0,
				Height: 0,
			},
		)

		require.Nil(t, img)
		require.Equal(t, ErrResourceUnavailable, err)
	})

	t.Run("unsupported image mime type", func(t *testing.T) {
		testFile, _ := os.Open("../../../tests/data/_gopher_original_1024x504.png")

		client := mocks.NewMockClient(ctrl)
		client.
			EXPECT().
			Get(gomock.Any()).
			Return(&http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(testFile),
			}, nil)

		img, err := NewHTTPDownloader(client).Download(
			"http://yandex.ru/test.jpg",
			dto.ImageDimensions{
				Width:  0,
				Height: 0,
			},
		)

		require.Nil(t, img)
		require.Error(t, ErrInvalidJpeg, err)
	})

	t.Run("invalid dimensions", func(t *testing.T) {
		testFile, _ := os.Open("../../../tests/data/_gopher_original_1024x504.jpg")

		client := mocks.NewMockClient(ctrl)
		client.
			EXPECT().
			Get(gomock.Any()).
			Return(&http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(testFile),
			}, nil)

		img, err := NewHTTPDownloader(client).Download(
			"http://yandex.ru/test.jpg",
			dto.ImageDimensions{
				Width:  10000,
				Height: 20000,
			},
		)

		require.Nil(t, img)
		require.Error(t, ErrSourceHasWrongDimensions, err)
	})

	t.Run("response should be valid", func(t *testing.T) {
		testFile, _ := os.Open("../../../tests/data/_gopher_original_1024x504.jpg")

		client := mocks.NewMockClient(ctrl)
		client.
			EXPECT().
			Get(gomock.Any()).
			Return(&http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(testFile),
			}, nil)

		img, err := NewHTTPDownloader(client).Download(
			"http://yandex.ru/test.jpg",
			dto.ImageDimensions{
				Width:  200,
				Height: 200,
			},
		)

		require.NotNil(t, img)
		require.Nil(t, err)
	})
}

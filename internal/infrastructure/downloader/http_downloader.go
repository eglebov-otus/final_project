package downloader

import (
	"errors"
	"github.com/disintegration/imaging"
	"go.uber.org/zap"
	"image"
	"image-previewer/internal/domain/valueObjects"
	"image/jpeg"
	"net/http"
)

var ErrResourceUnavailable = errors.New("image resource unavailable")
var ErrInvalidJpeg = errors.New("image should have correct jpeg struct")
var ErrSourceHasWrongDimensions = errors.New("source image has wrong dimensions")

type Client interface {
	Get(url string) (resp *http.Response, err error)
}

type httpClient struct {
}

func (c *httpClient) Get(url string) (resp *http.Response, err error) {
	return http.Get(url)
}

func NewHttpClient() *httpClient {
	return &httpClient{}
}

type httpDownloader struct {
	client Client
}

func (d *httpDownloader) Download(url string, dim valueObjects.ImageDimensions) (image.Image, error) {
	resp, err := d.client.Get("http://" + url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrResourceUnavailable
	}

	img, err := jpeg.Decode(resp.Body)

	if err != nil {
		return nil, ErrInvalidJpeg
	}

	bounds := img.Bounds()

	if bounds.Dx() < dim.Width || bounds.Dy() < dim.Height {
		return nil, ErrSourceHasWrongDimensions
	}

	zap.S().Debugf("crop downloaded image %d x %d", dim.Width, dim.Height)

	return imaging.Fill(img, dim.Width, dim.Height, imaging.Center, imaging.Lanczos), nil
}

func NewHttpDownloader(c Client) *httpDownloader {
	return &httpDownloader{
		client: c,
	}
}



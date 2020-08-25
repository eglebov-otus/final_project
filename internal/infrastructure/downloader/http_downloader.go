package downloader

import (
	"errors"
	"image"
	"image-previewer/internal/domain/dto"
	"image/jpeg"
	"net/http"

	"github.com/disintegration/imaging"
	"go.uber.org/zap"
)

var (
	ErrResourceUnavailable      = errors.New("image resource unavailable")
	ErrInvalidJpeg              = errors.New("image should have correct jpeg struct")
	ErrSourceHasWrongDimensions = errors.New("source image has wrong dimensions")
)

type Client interface {
	Get(url string) (resp *http.Response, err error)
}

type HTTPClient struct {
}

func (c *HTTPClient) Get(url string) (resp *http.Response, err error) {
	//nolint
	return http.Get(url)
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{}
}

type HTTPDownloader struct {
	client Client
}

func (d *HTTPDownloader) Download(url string, dim dto.ImageDimensions) (image.Image, error) {
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

func NewHTTPDownloader(c Client) *HTTPDownloader {
	return &HTTPDownloader{
		client: c,
	}
}

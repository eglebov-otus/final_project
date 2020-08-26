package downloader

import (
	"errors"
	"image"
	"image-previewer/internal/domain"
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
	Get(rawURL string, headers domain.RequestHeaders) (resp *http.Response, err error)
}

type HTTPClient struct {
}

func (c *HTTPClient) Get(rawURL string, headers domain.RequestHeaders) (resp *http.Response, err error) {
	client := http.Client{}

	//nolint:noctx
	req, err := http.NewRequest(http.MethodGet, "http://"+rawURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header(headers)

	return client.Do(req)
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{}
}

type HTTPDownloader struct {
	client Client
}

func (d *HTTPDownloader) Download(url string, dim dto.ImageDimensions, headers domain.RequestHeaders) (image.Image, error) {
	resp, err := d.client.Get(url, headers)
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

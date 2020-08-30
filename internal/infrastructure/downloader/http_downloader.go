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
	ErrResourceUnavailable = errors.New("image resource unavailable")
	ErrInvalidJpeg         = errors.New("image should have correct jpeg struct")
)

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

	zap.S().Debugf("resizing downloaded image %d x %d", dim.Width, dim.Height)

	return imaging.Fill(img, dim.Width, dim.Height, imaging.Center, imaging.Lanczos), nil
}

func NewHTTPDownloader(c Client) *HTTPDownloader {
	return &HTTPDownloader{
		client: c,
	}
}

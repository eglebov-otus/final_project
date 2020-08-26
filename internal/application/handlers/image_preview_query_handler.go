package handlers

import (
	"errors"
	"image"
	"image-previewer/internal/application/queries"
	"image-previewer/internal/domain"
	"net/url"

	"go.uber.org/zap"
)

var (
	ErrInvalidWidth  = errors.New("width should be greater than 0")
	ErrInvalidHeight = errors.New("height should be greater than 0")
	ErrEmptyURL      = errors.New("url should not be empty")
	ErrInvalidURL    = errors.New("url should be valid")
	ErrNotFound      = errors.New("img not found")
)

type ImagePreviewQueryHandler struct {
	previewRepository domain.PreviewRepository
	downloader        domain.Downloader
	idResolver        domain.ImageIDResolver
}

func (h *ImagePreviewQueryHandler) Handle(q queries.ImagePreviewQuery) (image.Image, error) {
	if err := h.checkQuery(q); err != nil {
		return nil, err
	}

	imageID := h.idResolver.ResolveImageID(q.URL, q.Dimensions)

	zap.S().Debugf("started processing image %s", string(imageID))

	img, err := h.previewRepository.FindOne(imageID)

	if err == ErrNotFound {
		zap.S().Debug("not found in cache, downloading")

		img, err = h.downloader.Download(q.URL, q.Dimensions, q.Headers)

		if err != nil {
			return nil, err
		}

		zap.S().Debug("adding to repository")

		_, err = h.previewRepository.Add(imageID, img)

		if err != nil {
			return nil, err
		}

		return img, nil
	}

	zap.S().Debug("using image from cache")

	return img, err
}

func (h *ImagePreviewQueryHandler) checkQuery(q queries.ImagePreviewQuery) error {
	if q.Dimensions.Width < 1 {
		return ErrInvalidWidth
	}

	if q.Dimensions.Height < 1 {
		return ErrInvalidHeight
	}

	if q.URL == "" {
		return ErrEmptyURL
	}

	if _, err := url.Parse(q.URL); err != nil {
		return ErrInvalidURL
	}

	return nil
}

func NewImagePreviewQueryHandler(
	rep domain.PreviewRepository,
	downloader domain.Downloader,
	resolver domain.ImageIDResolver,
) *ImagePreviewQueryHandler {
	return &ImagePreviewQueryHandler{
		previewRepository: rep,
		downloader:        downloader,
		idResolver:        resolver,
	}
}

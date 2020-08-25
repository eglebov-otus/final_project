package domain

import (
	"image"
	"image-previewer/internal/domain/dto"
)

type Downloader interface {
	Download(url string, dimensions dto.ImageDimensions) (image.Image, error)
}

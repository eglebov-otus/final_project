package domain

import (
	"image"
	"image-previewer/internal/domain/valueObjects"
)

type Downloader interface {
	Download(url string, dimensions valueObjects.ImageDimensions) (image.Image, error)
}

package domain

import (
	"image"
	"image-previewer/internal/domain/dto"
)

type RequestHeaders map[string][]string

type Downloader interface {
	Download(url string, dimensions dto.ImageDimensions, headers RequestHeaders) (image.Image, error)
}

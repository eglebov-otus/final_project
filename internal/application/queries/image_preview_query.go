package queries

import (
	"image-previewer/internal/domain"
	"image-previewer/internal/domain/dto"
)

type ImagePreviewQuery struct {
	URL        string
	Headers    domain.RequestHeaders
	Dimensions dto.ImageDimensions
}

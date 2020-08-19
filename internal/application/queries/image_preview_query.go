package queries

import "image-previewer/internal/domain/valueObjects"

type ImagePreviewQuery struct {
	Url string
	Dimensions valueObjects.ImageDimensions
}
package queries

import "image-previewer/internal/domain/dto"

type RequestHeaders map[string][]string

type ImagePreviewQuery struct {
	URL        string
	Headers    RequestHeaders
	Dimensions dto.ImageDimensions
}

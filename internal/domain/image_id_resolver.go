package domain

import "image-previewer/internal/domain/dto"

type ImageID string

type ImageIDResolver interface {
	ResolveImageID(url string, dim dto.ImageDimensions) ImageID
}

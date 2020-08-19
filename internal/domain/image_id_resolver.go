package domain

import "image-previewer/internal/domain/valueObjects"

type ImageId string

type ImageIdResolver interface {
	ResolveImageId(url string, dim valueObjects.ImageDimensions) ImageId
}

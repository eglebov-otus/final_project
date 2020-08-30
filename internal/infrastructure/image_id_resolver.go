package infrastructure

import (
	"fmt"
	"hash/fnv"
	"image-previewer/internal/domain"
	"image-previewer/internal/domain/dto"
)

type ImageIDResolver struct {
}

func (r *ImageIDResolver) ResolveImageID(url string, dim dto.ImageDimensions) domain.ImageID {
	h := fnv.New32a()
	_, _ = h.Write([]byte(url))

	imageID := fmt.Sprintf("%d_%dx%d", int(h.Sum32()), dim.Width, dim.Height)

	return domain.ImageID(imageID)
}

func NewImageIDResolver() *ImageIDResolver {
	return &ImageIDResolver{}
}

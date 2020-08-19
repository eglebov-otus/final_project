package infrastructure

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image-previewer/internal/domain"
	"image-previewer/internal/domain/valueObjects"
)

type imageIdResolver struct {
}

func (r *imageIdResolver) ResolveImageId(url string, dim valueObjects.ImageDimensions) domain.ImageId  {
	hash := md5.Sum([]byte(url))
	imageId := fmt.Sprintf("%s_%dx%d", hex.EncodeToString(hash[:]), dim.Width, dim.Height)

	return domain.ImageId(imageId)
}

func NewImageIdResolver() *imageIdResolver {
	return &imageIdResolver{}
}
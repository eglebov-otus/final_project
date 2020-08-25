package infrastructure

//nolint:gosec
import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image-previewer/internal/domain"
	"image-previewer/internal/domain/dto"
)

type ImageIDResolver struct {
}

func (r *ImageIDResolver) ResolveImageID(url string, dim dto.ImageDimensions) domain.ImageID {
	//nolint:gosec
	hash := md5.Sum([]byte(url))
	imageID := fmt.Sprintf("%s_%dx%d", hex.EncodeToString(hash[:]), dim.Width, dim.Height)

	return domain.ImageID(imageID)
}

func NewImageIDResolver() *ImageIDResolver {
	return &ImageIDResolver{}
}

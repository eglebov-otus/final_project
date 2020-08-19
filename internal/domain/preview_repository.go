package domain

import "image"

type PreviewRepository interface {
	FindOne(id ImageId) (image.Image, error)
	Add(id ImageId, img image.Image) (bool, error)
}

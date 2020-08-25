package domain

import "image"

type PreviewRepository interface {
	FindOne(id ImageID) (image.Image, error)
	Add(id ImageID, img image.Image) (bool, error)
}

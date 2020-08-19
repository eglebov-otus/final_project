package infrastructure

import (
	"github.com/stretchr/testify/require"
	"image-previewer/internal/domain"
	"image-previewer/internal/domain/valueObjects"
	"testing"
)

func TestImageIdResolver_ResolveImageId(t *testing.T) {
	t.Run("resolve id should return valid string", func(t *testing.T) {
		actualId := NewImageIdResolver().ResolveImageId(
			"http://ya.ru/test.jpg",
			valueObjects.ImageDimensions{
				Width: 100,
				Height: 500,
			},
		)

		require.Equal(t, domain.ImageId("9508dfb97b74094e1b8134e15469fc0e_100x500"), actualId)
	})
}
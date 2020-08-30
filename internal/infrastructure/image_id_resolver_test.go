package infrastructure

import (
	"image-previewer/internal/domain"
	"image-previewer/internal/domain/dto"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImageIdResolver_ResolveImageId(t *testing.T) {
	t.Run("resolve id should return valid string", func(t *testing.T) {
		actualID := NewImageIDResolver().ResolveImageID(
			"http://ya.ru/test.jpg",
			dto.ImageDimensions{
				Width:  100,
				Height: 500,
			},
		)

		require.Equal(t, domain.ImageID("1864972448_100x500"), actualID)
	})
}

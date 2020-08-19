package handlers

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"image"
	"image-previewer/internal/application/queries"
	"image-previewer/internal/domain"
	"image-previewer/internal/domain/valueObjects"
	"image-previewer/tests/mocks"
	"image/jpeg"
	"os"
	"testing"
)

//go:generate mockgen -destination=../../../tests/mocks/mock_preview_repository.go -package=mocks image-previewer/internal/domain PreviewRepository
//go:generate mockgen -destination=../../../tests/mocks/mock_downloader.go -package=mocks image-previewer/internal/domain Downloader
//go:generate mockgen -destination=../../../tests/mocks/mock_id_resolver.go -package=mocks image-previewer/internal/domain ImageIdResolver
func TestImagePreviewQueryHandler_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Run("invalid query width and height", func(t *testing.T) {
		rep := mocks.NewMockPreviewRepository(ctrl)
		idResolver := mocks.NewMockImageIdResolver(ctrl)
		downloader := mocks.NewMockDownloader(ctrl)

		handler := NewImagePreviewQueryHandler(rep, downloader, idResolver)

		img, err := handler.Handle(queries.ImagePreviewQuery{
			Url: "http://ya.ru",
			Dimensions: valueObjects.ImageDimensions{
				Width: 0,
				Height: 200,
			},
		})

		require.Nil(t, img)
		require.Equal(t, err, ErrInvalidWidth)

		img, err = handler.Handle(queries.ImagePreviewQuery{
			Url: "http://ya.ru",
			Dimensions: valueObjects.ImageDimensions{
				Width: 100,
				Height: 0,
			},
		})

		require.Nil(t, img)
		require.Equal(t, err, ErrInvalidHeight)
	})

	t.Run("invalid query url", func(t *testing.T) {
		rep := mocks.NewMockPreviewRepository(ctrl)
		idResolver := mocks.NewMockImageIdResolver(ctrl)
		downloader := mocks.NewMockDownloader(ctrl)

		handler := NewImagePreviewQueryHandler(rep, downloader, idResolver)

		img, err := handler.Handle(queries.ImagePreviewQuery{
			Url: "",
			Dimensions: valueObjects.ImageDimensions{
				Width: 100,
				Height: 100,
			},
		})

		require.Nil(t, img)
		require.Equal(t, err, ErrEmptyUrl)
	})

	t.Run("image found in repository", func(t *testing.T) {
		rep := mocks.NewMockPreviewRepository(ctrl)
		rep.
			EXPECT().
			FindOne(gomock.Any()).
			Return(fakedImg(), nil)
		rep.
			EXPECT().
			Add(gomock.Any(), gomock.Any()).
			Times(0)

		idResolver := mocks.NewMockImageIdResolver(ctrl)
		idResolver.
			EXPECT().
			ResolveImageId(gomock.Any(), gomock.Any()).
			Return(domain.ImageId("test_id"))

		downloader := mocks.NewMockDownloader(ctrl)
		downloader.
			EXPECT().
			Download(gomock.Any(), gomock.Any()).
			Times(0)

		handler := NewImagePreviewQueryHandler(rep, downloader, idResolver)

		img, err := handler.Handle(queries.ImagePreviewQuery{
			Url: "http://ya.ru",
			Dimensions: valueObjects.ImageDimensions{
				Width: 100,
				Height: 200,
			},
		})

		require.Nil(t, err)
		require.NotNil(t, img)
	})

	t.Run("image not found in repository", func(t *testing.T) {
		rep := mocks.NewMockPreviewRepository(ctrl)
		rep.
			EXPECT().
			FindOne(gomock.Any()).
			Return(nil, ErrNotFound)
		rep.
			EXPECT().
			Add(gomock.Any(), gomock.Any()).
			Return(true, nil).
			Times(1)

		idResolver := mocks.NewMockImageIdResolver(ctrl)
		idResolver.
			EXPECT().
			ResolveImageId(gomock.Any(), gomock.Any()).
			Return(domain.ImageId("test_id"))

		actualImg := fakedImg()

		downloader := mocks.NewMockDownloader(ctrl)
		downloader.
			EXPECT().
			Download(gomock.Any(), gomock.Any()).
			Return(actualImg, nil).
			Times(1)

		handler := NewImagePreviewQueryHandler(rep, downloader, idResolver)

		img, err := handler.Handle(queries.ImagePreviewQuery{
			Url: "http://ya.ru",
			Dimensions: valueObjects.ImageDimensions{
				Width: 100,
				Height: 200,
			},
		})

		require.Nil(t, err)
		require.NotNil(t, img)
		require.Same(t, img, actualImg)
	})

}

func fakedImg() image.Image {
	f, _ := os.Open("../../../tests/data/_gopher_500x500.jpg")
	img, _ := jpeg.Decode(f)
	f.Close()

	return img
}
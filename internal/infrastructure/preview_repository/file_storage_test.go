package preview_repository

import (
	"github.com/stretchr/testify/require"
	"image"
	"image-previewer/internal/application/handlers"
	"image-previewer/internal/domain"
	"image/jpeg"
	"io/ioutil"
	"os"
	"testing"
)

var cacheDir = "../../../tests/data/cache/"

func TestFileStorage_Add(t *testing.T) {

	t.Run("valid response status", func(t *testing.T) {
		defer cleanUp(cacheDir)

		s := NewFileStorage(cacheDir, 5)

		require.Equal(t, 0, s.Len())

		wasInCache, err := s.Add(domain.ImageId("test1.jpg"), fakedImg())
		require.False(t, wasInCache)
		require.Nil(t, err)
		_, err = os.Open(cacheDir + "test1.jpg")
		require.Nil(t, err)

		wasInCache, err = s.Add(domain.ImageId("test1.jpg"), fakedImg())
		require.True(t, wasInCache)
		require.Nil(t, err)
		_, err = os.Open(cacheDir + "test1.jpg")
		require.Nil(t, err)

		require.Equal(t, 1, s.Len())

		wasInCache, err = s.Add(domain.ImageId("test2.jpg"), fakedImg())
		require.False(t, wasInCache)
		require.Nil(t, err)
		_, err = os.Open(cacheDir + "test2.jpg")
		require.Nil(t, err)

		require.Equal(t, 2, s.Len())
	})

	t.Run("purge logic", func(t *testing.T) {
		defer cleanUp(cacheDir)

		s := NewFileStorage(cacheDir, 3)
		_, _ = s.Add(domain.ImageId("test1.jpg"), fakedImg())
		require.Equal(t, 1, s.Len())
		_, _ = s.Add(domain.ImageId("test2.jpg"), fakedImg())
		require.Equal(t, 2, s.Len())
		_, _ = s.Add(domain.ImageId("test3.jpg"), fakedImg())
		require.Equal(t, 3, s.Len())
		_, _ = s.Add(domain.ImageId("test4.jpg"), fakedImg())
		require.Equal(t, 3, s.Len())
		_, _ = s.Add(domain.ImageId("test5.jpg"), fakedImg())
		require.Equal(t, 3, s.Len())

		_, err := os.Open(cacheDir + "test1.jpg")
		require.NotNil(t, err)
		_, err = os.Open(cacheDir + "test2.jpg")
		require.NotNil(t, err)
	})
}

func TestFileStorage_FindOne(t *testing.T) {
	defer cleanUp(cacheDir)

	t.Run("not found case", func(t *testing.T) {
		s := NewFileStorage(cacheDir, 5)

		img, err := s.FindOne(domain.ImageId("test500.jpg"))

		require.Nil(t, img)
		require.Equal(t, err, handlers.ErrNotFound)
	})

	t.Run("found case", func(t *testing.T) {
		s := NewFileStorage(cacheDir, 5)
		imageId := domain.ImageId("test500.jpg")
		_, _ = s.Add(imageId, fakedImg())

		img, err := s.FindOne(imageId)

		require.NotNil(t, img)
		require.Nil(t, err)
	})
}

func fakedImg() image.Image {
	f, _ := os.Open("../../../tests/data/_gopher_500x500.jpg")
	img, _ := jpeg.Decode(f)
	f.Close()

	return img
}

func cleanUp(cacheDir string) {
	files, _ := ioutil.ReadDir(cacheDir)

	for _, file := range files {
		if file.Name() != ".gitkeep" {
			_ = os.Remove(cacheDir + file.Name())
		}
	}
}
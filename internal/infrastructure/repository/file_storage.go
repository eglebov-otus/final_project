package repository

import (
	"container/list"
	"fmt"
	"image"
	"image-previewer/internal/application/handlers"
	"image-previewer/internal/domain"
	"image/jpeg"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
)

type FileStorage struct {
	cacheDir string
	capacity int
	cache    list.List
	items    map[domain.ImageID]*list.Element
	mux      sync.Mutex
}

func (r *FileStorage) FindOne(id domain.ImageID) (image.Image, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	element, exists := r.items[id]

	if !exists {
		return nil, handlers.ErrNotFound
	}

	if err := r.touchPreview(id); err != nil {
		return nil, err
	}

	r.cache.MoveToFront(element)

	img, err := r.loadPreview(id)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (r *FileStorage) Add(id domain.ImageID, img image.Image) (bool, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	element, exists := r.items[id]

	if exists {
		zap.S().Debugf("item exist in cache, moving to front")

		if err := r.touchPreview(id); err != nil {
			return true, err
		}

		r.cache.MoveToFront(element)
		element.Value = id

		return true, nil
	}

	if r.cache.Len() == r.capacity {
		zap.S().Debugf("cache capacity limit exceed, removing last item")

		lastItem := r.cache.Back()
		lastItemID := lastItem.Value.(domain.ImageID)

		if err := r.removePreview(lastItemID); err != nil {
			return false, err
		}

		r.cache.Remove(lastItem)
		delete(r.items, lastItemID)
	}

	zap.S().Debugf("new item, saving and pushing to front")

	if err := r.savePreview(id, img); err != nil {
		return false, err
	}

	element = r.cache.PushFront(id)
	r.items[id] = element

	return false, nil
}

func (r *FileStorage) Len() int {
	return r.cache.Len()
}

func (r *FileStorage) savePreview(id domain.ImageID, img image.Image) error {
	path := r.pathByID(id)

	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %s", path, err)
	}

	defer out.Close()

	if err := jpeg.Encode(out, img, nil); err != nil {
		return fmt.Errorf("failed to encode image %s: %s", path, err)
	}

	return nil
}

func (r *FileStorage) loadPreview(id domain.ImageID) (image.Image, error) {
	path := r.pathByID(id)

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %s", path, err)
	}

	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file %s: %s", path, err)
	}

	return img, nil
}

func (r *FileStorage) removePreview(id domain.ImageID) error {
	path := r.pathByID(id)

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to remove file %s: %s", path, err)
	}

	return nil
}

func (r *FileStorage) touchPreview(id domain.ImageID) error {
	path := r.pathByID(id)

	if err := os.Chtimes(path, time.Now(), time.Now()); err != nil {
		return fmt.Errorf("failed to touch file %s: %s", path, err)
	}

	return nil
}

func (r *FileStorage) pathByID(id domain.ImageID) string {
	return r.cacheDir + string(id)
}

func NewFileStorage(cacheDir string, capacity int) *FileStorage {
	return &FileStorage{
		cacheDir: cacheDir,
		capacity: capacity,
		items:    make(map[domain.ImageID]*list.Element),
	}
}

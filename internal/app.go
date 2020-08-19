package internal

import (
	"bytes"
	"errors"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"image-previewer/internal/application/handlers"
	"image-previewer/internal/application/queries"
	"image-previewer/internal/domain/valueObjects"
	"image-previewer/internal/infrastructure"
	"image-previewer/internal/infrastructure/downloader"
	"image-previewer/internal/infrastructure/preview_repository"
	"image/jpeg"
	"net/http"
	"strconv"
)

type app struct {
}

func NewApp() *app {
	return &app{}
}

func (app *app) Run() error {
	cacheDir := viper.GetString("app.preview_cache_dir")
	capacity := viper.GetInt("app.preview_cache_size")

	if capacity == 0 {
		return errors.New("invalid config: preview_cache_size should be set")
	}

	if cacheDir == "" {
		return errors.New("invalid config: preview_cache_dir should be set")
	}

	rep := preview_repository.NewFileStorage(cacheDir, capacity)
	idResolver := infrastructure.NewImageIdResolver()
	httpDownloader := downloader.NewHttpDownloader(downloader.NewHttpClient())

	queryHandler := handlers.NewImagePreviewQueryHandler(rep, httpDownloader, idResolver)

	router := mux.NewRouter()
	router.HandleFunc("/fill/{width}/{height}/{url:.*}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		width, err := strconv.Atoi(vars["width"])

		if err != nil {
			zap.S().Warnf("invalid width value", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		height, err := strconv.Atoi(vars["height"])

		if err != nil {
			zap.S().Warnf("invalid height value", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		img, err := queryHandler.Handle(queries.ImagePreviewQuery{
			Url: vars["url"],
			Dimensions: valueObjects.ImageDimensions{
				Width: width,
				Height: height,
			},
		})

		if err != nil {
			zap.S().Errorf("get preview query handle failed: %s", err)

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		buf := new(bytes.Buffer)

		if err := jpeg.Encode(buf, img, nil); err != nil {
			zap.S().Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(buf.Bytes()); err != nil {
			zap.S().Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(buf.Bytes())))
		w.WriteHeader(http.StatusOK)
	})

	http.Handle("/", router)

	zap.S().Info("Application started")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		return err
	}

	return nil
}

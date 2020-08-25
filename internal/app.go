package internal

import (
	"errors"
	"image-previewer/internal/application/handlers"
	"image-previewer/internal/infrastructure"
	"image-previewer/internal/infrastructure/downloader"
	"image-previewer/internal/infrastructure/repository"
	"image-previewer/internal/interfaces/http/controllers"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (app *App) Run() error {
	cacheDir := viper.GetString("app.preview_cache_dir")
	capacity := viper.GetInt("app.preview_cache_size")

	if capacity == 0 {
		return errors.New("invalid config: preview_cache_size should be set")
	}

	if cacheDir == "" {
		return errors.New("invalid config: preview_cache_dir should be set")
	}

	rep := repository.NewFileStorage(cacheDir, capacity)
	idResolver := infrastructure.NewImageIDResolver()
	httpDownloader := downloader.NewHTTPDownloader(downloader.NewHTTPClient())
	queryHandler := handlers.NewImagePreviewQueryHandler(rep, httpDownloader, idResolver)
	controller := controllers.NewImagePreviewController(queryHandler)

	router := mux.NewRouter()
	router.HandleFunc("/fill/{width}/{height}/{url:.*}", controller.ActionGet)

	http.Handle("/", router)

	zap.S().Info("Application started")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		return err
	}

	return nil
}

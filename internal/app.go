package internal

import (
	"context"
	"errors"
	"image-previewer/internal/application/handlers"
	"image-previewer/internal/infrastructure"
	"image-previewer/internal/infrastructure/downloader"
	"image-previewer/internal/infrastructure/repository"
	"image-previewer/internal/interfaces/http/controllers"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-signalCh
		cancel()
	}()

	if err := serve(ctx, cacheDir, capacity); err != nil {
		zap.S().Fatalf("failed to serve: %s", err)

		return err
	}

	return nil
}

func serve(ctx context.Context, cacheDir string, capacity int) (err error) {
	rep := repository.NewFileStorage(cacheDir, capacity)
	idResolver := infrastructure.NewImageIDResolver()
	httpDownloader := downloader.NewHTTPDownloader(downloader.NewHTTPClient(&http.Client{}))
	queryHandler := handlers.NewImagePreviewQueryHandler(rep, httpDownloader, idResolver)
	controller := controllers.NewImagePreviewController(queryHandler)

	router := mux.NewRouter()
	router.HandleFunc("/fill/{width}/{height}/{url:.*}", controller.ActionGet)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.S().Fatalf("Failed to run server: %s", err)
		}
	}()

	zap.S().Info("server started")

	<-ctx.Done()

	zap.S().Info("stopping server")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		zap.S().Fatalf("shutdown failed: %s", err)
	}

	zap.S().Info("server stopped")

	if err == http.ErrServerClosed {
		err = nil
	}

	return
}

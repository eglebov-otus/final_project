package controllers

import (
	"bytes"
	"image-previewer/internal/application/handlers"
	"image-previewer/internal/application/queries"
	"image-previewer/internal/domain"
	"image-previewer/internal/domain/dto"
	"image/jpeg"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ImagePreviewController struct {
	handler *handlers.ImagePreviewQueryHandler
}

func (c *ImagePreviewController) ActionGet(w http.ResponseWriter, r *http.Request) {
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

	img, err := c.handler.Handle(queries.ImagePreviewQuery{
		URL:     vars["url"],
		Headers: domain.RequestHeaders(r.Header),
		Dimensions: dto.ImageDimensions{
			Width:  width,
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
}

func NewImagePreviewController(h *handlers.ImagePreviewQueryHandler) *ImagePreviewController {
	return &ImagePreviewController{
		handler: h,
	}
}

package handler

import (
	"backend/internal/handler/dto/mapper"
	"backend/internal/service"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

type TrackHandler struct {
	svc *service.TrackService
}

func NewTrackHdl(svc *service.TrackService) *TrackHandler {
	return &TrackHandler{svc: svc}
}

func (h *TrackHandler) Upload(c echo.Context) error {
	file, err := c.FormFile("track")
	if err != nil {
		response := mapper.BadRequestErrors("Файл не предоставлен")
		return c.JSON(http.StatusBadRequest, response)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{
		".mp3":  true,
		".wav":  true,
		".ogg":  true,
		".flac": true,
		".aac":  true,
		".m4a":  true,
	}
	if !allowed[ext] {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Not .mp3, .wav, .ogg, .flac, .aac, .m4a"})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Не удалось открыть файл"})
	}
	defer src.Close()

	track, err := h.svc.Upload(src, file.Filename)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, track)
}

func (h *TrackHandler) Search(c echo.Context) error {
	q := c.QueryParam("q")
	tracks, err := h.svc.Search(q)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, tracks)
}

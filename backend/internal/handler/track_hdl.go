package handler

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"localtunes/internal/service"
)

type TrackHdl struct {
	svc *service.TrackSvc
}

func NewTrackHdl(svc *service.TrackSvc) *TrackHdl {
	return &TrackHdl{svc: svc}
}

func (h *TrackHdl) Upload(c echo.Context) error {
	file, err := c.FormFile("track")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Файл не предоставлен"})
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

func (h *TrackHdl) Search(c echo.Context) error {
	q := c.QueryParam("q")
	tracks, err := h.svc.Search(q)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, tracks)
}
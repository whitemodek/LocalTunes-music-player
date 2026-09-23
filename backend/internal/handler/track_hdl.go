package handler

import (
	"backend/internal/handler/dto/mapper"
	"backend/internal/service"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

type TrackHandler struct {
	svc        *service.TrackService
	uploadsDir string
}

func NewTrackHdl(svc *service.TrackService, uploadsDir string) *TrackHandler {
	return &TrackHandler{svc: svc, uploadsDir: uploadsDir}
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
		return c.JSON(http.StatusBadRequest, mapper.BadRequestErrors("Not .mp3, .wav, .ogg, .flac, .aac, .m4a"))
	}

	src, err := file.Open()
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, mapper.InternalServerError())
	}
	defer src.Close()

	track, err := h.svc.Upload(src, file.Filename)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, mapper.InternalServerError())
	}

	return c.JSON(http.StatusOK, track)
}

func (h *TrackHandler) Search(c echo.Context) error {
	q := c.QueryParam("q")
	tracks, err := h.svc.Search(q)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, mapper.InternalServerError())
	}
	return c.JSON(http.StatusOK, tracks)
}

func (h *TrackHandler) StreamDirect(c echo.Context) error {
	trackID := c.Param("track_id")
	name := filepath.Base(trackID)
	if name != trackID || name == "." || name == ".." {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid track name")
	}
	filePath := filepath.Join(h.uploadsDir, name)

	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return echo.NewHTTPError(http.StatusNotFound, "track not found")
		}
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return c.File(filePath)
}

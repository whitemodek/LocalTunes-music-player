package main

import (
	"backend/internal/handler"
	repo "backend/internal/repository"
	"backend/internal/service"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))

	trackRepo, err := repo.NewTrackRepo("./music.db")
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer trackRepo.Close()

	trackSvc := service.NewTrackSvc(trackRepo)
	trackHdl := handler.NewTrackHdl(trackSvc)

	api := e.Group("/api")
	{
		api.POST("/upload", trackHdl.Upload)
		api.GET("/search", trackHdl.Search)
	}

	e.GET("/stream-direct/:track_id", func(c echo.Context) error {
		trackID := c.Param("track_id")
		filePath := filepath.Join("uploads", trackID)

		if _, err := os.Stat(filePath); err != nil {
			if os.IsNotExist(err) {
				return echo.NewHTTPError(http.StatusNotFound, "track not found")
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.File(filePath)
	})

	e.Static("/api/stream", "uploads")
	e.Static("/static", "static")

	log.Println("Сервер запущен на :8080")
	e.Start(":8080")
}

package main

import (
	"backend/internal/config"
	"backend/internal/handler"
	repo "backend/internal/repository"
	"backend/internal/service"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.BodyLimit(cfg.BodyLimit))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))

	trackRepo, err := repo.NewTrackRepo(cfg.DBPath)
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer trackRepo.Close()

	trackSvc := service.NewTrackSvc(trackRepo, cfg.UploadsDir)
	trackHdl := handler.NewTrackHdl(trackSvc, cfg.UploadsDir)

	api := e.Group("/api")
	{
		api.POST("/upload", trackHdl.Upload)
		api.GET("/search", trackHdl.Search)
	}

	e.GET("/stream-direct/:track_id", trackHdl.StreamDirect)
	e.Static("/api/stream", cfg.UploadsDir)
	e.Static("/static", "static")

	log.Printf("Сервер запущен на %s", cfg.Port)
	e.Start(cfg.Port)
}

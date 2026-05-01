package main

import (
	"log"
	"net/http"
	"localtunes/internal/handler"
	"localtunes/internal/repo"
	"localtunes/internal/service"
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

	e.Static("/api/stream", "uploads")
	e.Static("/static", "static")

	log.Println("Сервер запущен на :8080")
	e.Start(":8080")
}

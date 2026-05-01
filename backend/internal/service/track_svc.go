package service

import (
	"localtunes/internal/models"
	"localtunes/internal/repository"
)

func ListTracks() []models.Track {
	return repository.GetAllTracks()
}

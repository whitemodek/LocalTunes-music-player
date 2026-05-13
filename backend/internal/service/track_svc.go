package service

import (
	"backend/internal/models"
	repo "backend/internal/repository"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type TrackService struct {
	repo *repo.TrackRepo
}

func NewTrackSvc(r *repo.TrackRepo) *TrackService {
	return &TrackService{repo: r}
}

func (s *TrackService) Upload(src io.Reader, origName string) (*models.Track, error) {
	os.MkdirAll("uploads", 0755)

	ext := filepath.Ext(origName)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join("uploads", fileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()
	io.Copy(dst, src)

	track := &models.Track{
		FileName:  fileName,
		StreamURL: "/api/stream/" + fileName,
	}

	f, err := os.Open(filePath)
	if err == nil {
		defer f.Close()
	}

	if track.Title == "" {
		track.Title = origName
	}
	if track.Artist == "" {
		track.Artist = "Unknown Artist"
	}

	if err := s.repo.Save(track); err != nil {
		return nil, err
	}

	return track, nil
}

func (s *TrackService) Search(q string) ([]models.Track, error) {
	return s.repo.Search(q)
}

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
	repo       *repo.TrackRepo
	uploadsDir string
}

func NewTrackSvc(r *repo.TrackRepo, uploadsDir string) *TrackService {
	return &TrackService{repo: r, uploadsDir: uploadsDir}
}

func (s *TrackService) Upload(src io.Reader, origName string) (*models.Track, error) {
	if err := os.MkdirAll(s.uploadsDir, 0755); err != nil {
		return nil, err
	}

	ext := filepath.Ext(origName)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(s.uploadsDir, fileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	track := &models.Track{
		FileName:  fileName,
		StreamURL: "/api/stream/" + fileName,
		Title:     origName,
		Artist:    "Unknown Artist",
	}

	if err := s.repo.Save(track); err != nil {
		return nil, err
	}

	return track, nil
}

func (s *TrackService) Search(q string) ([]models.Track, error) {
	return s.repo.Search(q)
}

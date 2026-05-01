package service

import (
	"backend/internal/models"
	repo "backend/internal/repository"
	repository "backend/internal/repository"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type TrackSvc struct {
	repo *repository.TrackRepo
}

func NewTrackSvc(repo *repo.TrackRepo) *TrackSvc {
	return &TrackSvc{repo: repo}
}

func (s *TrackSvc) Upload(src io.Reader, origName string) (*models.Track, error) {
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
		CoverURL:  "/static/default.jpg", // Заглушка обложки
		StreamURL: "/api/stream/" + fileName,
	}

	f, err := os.Open(filePath)
	if err == nil {
		defer f.Close()
		// m, err := tag.ReadFrom(f)
		// if err == nil && m != nil {
		// 	track.Title = m.Title()
		// 	track.Artist = m.Artist()
		// 	track.Album = m.Album()
		// }
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

func (s *TrackSvc) Search(q string) ([]models.Track, error) {
	return s.repo.Search(q)
}

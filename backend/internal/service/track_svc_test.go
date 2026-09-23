package service

import (
	repo "backend/internal/repository"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newTestSvc(t *testing.T) (*TrackService, string) {
	t.Helper()
	r, err := repo.NewTrackRepo(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("NewTrackRepo: %v", err)
	}
	t.Cleanup(r.Close)

	uploadsDir := filepath.Join(t.TempDir(), "uploads")
	return NewTrackSvc(r, uploadsDir), uploadsDir
}

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) {
	return 0, errors.New("read failed")
}

func TestUploadWritesFileAndSavesTrack(t *testing.T) {
	svc, uploadsDir := newTestSvc(t)

	track, err := svc.Upload(strings.NewReader("audio-data"), "song.mp3")
	if err != nil {
		t.Fatalf("Upload: %v", err)
	}
	if track.ID == 0 {
		t.Error("track ID not set")
	}
	if track.Title != "song.mp3" {
		t.Errorf("Title = %q, want %q", track.Title, "song.mp3")
	}
	if track.Artist != "Unknown Artist" {
		t.Errorf("Artist = %q, want %q", track.Artist, "Unknown Artist")
	}
	if track.FileName == "" {
		t.Error("FileName is empty")
	}
	if track.StreamURL != "/api/stream/"+track.FileName {
		t.Errorf("StreamURL = %q, want %q", track.StreamURL, "/api/stream/"+track.FileName)
	}

	filePath := filepath.Join(uploadsDir, track.FileName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("reading uploaded file: %v", err)
	}
	if string(data) != "audio-data" {
		t.Errorf("file content = %q, want %q", data, "audio-data")
	}
}

func TestUploadPropagatesReadError(t *testing.T) {
	svc, _ := newTestSvc(t)

	if _, err := svc.Upload(failingReader{}, "song.mp3"); err == nil {
		t.Fatal("Upload with failing reader: want error, got nil")
	}
}

func TestSearchDelegatesToRepo(t *testing.T) {
	svc, _ := newTestSvc(t)

	if _, err := svc.Upload(strings.NewReader("data"), "findme.mp3"); err != nil {
		t.Fatalf("Upload: %v", err)
	}

	tracks, err := svc.Search("findme")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(tracks) != 1 {
		t.Fatalf("Search returned %d tracks, want 1", len(tracks))
	}
}

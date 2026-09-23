package repo

import (
	"backend/internal/models"
	"path/filepath"
	"testing"
)

func newTestRepo(t *testing.T) *TrackRepo {
	t.Helper()
	r, err := NewTrackRepo(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("NewTrackRepo: %v", err)
	}
	t.Cleanup(r.Close)
	return r
}

func TestSaveSetsID(t *testing.T) {
	r := newTestRepo(t)

	track := &models.Track{
		Title:     "Song",
		Artist:    "Artist",
		Album:     "Album",
		FileName:  "1.mp3",
		StreamURL: "/api/stream/1.mp3",
	}
	if err := r.Save(track); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if track.ID == 0 {
		t.Fatal("Save did not set track ID")
	}
}

func TestSaveDuplicateUpdatesAndKeepsID(t *testing.T) {
	r := newTestRepo(t)

	first := &models.Track{
		Title:    "Original",
		FileName: "dup.mp3",
	}
	if err := r.Save(first); err != nil {
		t.Fatalf("Save first: %v", err)
	}

	second := &models.Track{
		Title:    "Updated",
		FileName: "dup.mp3",
	}
	if err := r.Save(second); err != nil {
		t.Fatalf("Save duplicate: %v", err)
	}
	if second.ID != first.ID {
		t.Fatalf("duplicate save ID = %d, want %d", second.ID, first.ID)
	}

	tracks, err := r.Search("Updated")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(tracks) != 1 {
		t.Fatalf("Search returned %d tracks, want 1", len(tracks))
	}
	if tracks[0].Title != "Updated" {
		t.Fatalf("title = %q, want %q", tracks[0].Title, "Updated")
	}
}

func TestSearchByAlbum(t *testing.T) {
	r := newTestRepo(t)

	track := &models.Track{
		Title:    "Come Together",
		Artist:   "The Beatles",
		Album:    "Abbey Road",
		FileName: "ct.mp3",
	}
	if err := r.Save(track); err != nil {
		t.Fatalf("Save: %v", err)
	}

	tracks, err := r.Search("Abbey")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(tracks) != 1 {
		t.Fatalf("Search by album returned %d tracks, want 1", len(tracks))
	}
	if tracks[0].Album != "Abbey Road" {
		t.Fatalf("album = %q, want %q", tracks[0].Album, "Abbey Road")
	}
}

func TestSearchByTitleAndArtist(t *testing.T) {
	r := newTestRepo(t)

	track := &models.Track{
		Title:    "Yesterday",
		Artist:   "The Beatles",
		Album:    "Help!",
		FileName: "y.mp3",
	}
	if err := r.Save(track); err != nil {
		t.Fatalf("Save: %v", err)
	}

	for _, q := range []string{"Yesterday", "Beatles"} {
		tracks, err := r.Search(q)
		if err != nil {
			t.Fatalf("Search(%q): %v", q, err)
		}
		if len(tracks) != 1 {
			t.Fatalf("Search(%q) returned %d tracks, want 1", q, len(tracks))
		}
	}
}

func TestSearchNoMatch(t *testing.T) {
	r := newTestRepo(t)

	tracks, err := r.Search("zzz-no-match")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(tracks) != 0 {
		t.Fatalf("Search returned %d tracks, want 0", len(tracks))
	}
}

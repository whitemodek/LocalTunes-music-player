package repo

import (
	"database/sql"
	"fmt"
	"localtunes/internal/models"
	"modernc.org/sqlite"
)

type TrackRepo struct {
	db *sql.DB
}

func NewTrackRepo(dbPath string) (*TrackRepo, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tracks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		artist TEXT,
		album TEXT,
		cover_url TEXT,
		file_name TEXT UNIQUE
	)`)
	if err != nil {
		return nil, err
	}

	return &TrackRepo{db: db}, nil
}

func (r *TrackRepo) Save(t *models.Track) error {
	_, err := r.db.Exec(`INSERT OR IGNORE INTO tracks (title, artist, album, cover_url, file_name)
		VALUES (?, ?, ?, ?, ?)`, t.Title, t.Artist, t.Album, t.CoverURL, t.FileName)
	return err
}

func (r *TrackRepo) Search(q string) ([]models.Track, error) {
	rows, err := r.db.Query(`SELECT id, title, artist, album, cover_url, file_name 
		FROM tracks 
		WHERE title LIKE ? OR artist LIKE ? 
		ORDER BY id DESC`,
		fmt.Sprintf("%%%s%%", q), fmt.Sprintf("%%%s%%", q))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []models.Track
	for rows.Next() {
		var t models.Track
		var fileName string
		if err := rows.Scan(&t.ID, &t.Title, &t.Artist, &t.Album, &t.CoverURL, &fileName); err != nil {
			return nil, err
		}
		t.StreamURL = "/api/stream/" + fileName
		tracks = append(tracks, t)
	}
	return tracks, nil
}

func (r *TrackRepo) Close() {
	r.db.Close()
}
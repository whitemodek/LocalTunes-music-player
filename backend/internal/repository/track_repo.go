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
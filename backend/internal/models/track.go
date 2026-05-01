package models

type Track struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	Album     string `json:"album"`
	CoverURL  string `json:"coverUrl"`
	StreamURL string `json:"streamUrl"`
	FileName  string `json:"-"`
}
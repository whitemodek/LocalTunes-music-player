package repository

import "localtunes/internal/models"

var sampleTracks = []models.Track{
	{ID: "1", Title: "Local Sunrise", Artist: "The Loop", Duration: 210},
	{ID: "2", Title: "Late Night Vinyl", Artist: "Mono Soul", Duration: 185},
}

func GetAllTracks() []models.Track {
	return sampleTracks
}

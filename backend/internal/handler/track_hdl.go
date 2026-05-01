package handler

import (
	"encoding/json"
	"net/http"

	"localtunes/internal/service"
)

func TrackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tracks := service.ListTracks()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tracks)
}

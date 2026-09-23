package handler

import (
	repo "backend/internal/repository"
	"backend/internal/service"
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/labstack/echo/v4"
)

func newTestHandler(t *testing.T) (*TrackHandler, string) {
	t.Helper()
	r, err := repo.NewTrackRepo(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("NewTrackRepo: %v", err)
	}
	t.Cleanup(r.Close)

	uploadsDir := filepath.Join(t.TempDir(), "uploads")
	svc := service.NewTrackSvc(r, uploadsDir)
	return NewTrackHdl(svc, uploadsDir), uploadsDir
}

func multipartTrack(t *testing.T, fileName string, content []byte) (*bytes.Buffer, string) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("track", fileName)
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	if _, err := fw.Write(content); err != nil {
		t.Fatalf("write part: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	return &buf, w.FormDataContentType()
}

func TestUploadMissingFile(t *testing.T) {
	h, _ := newTestHandler(t)
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/api/upload", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()

	if err := h.Upload(e.NewContext(req, rec)); err != nil {
		t.Fatalf("Upload: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestUploadBadExtension(t *testing.T) {
	h, _ := newTestHandler(t)
	e := echo.New()

	body, contentType := multipartTrack(t, "payload.exe", []byte("MZ"))
	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set(echo.HeaderContentType, contentType)
	rec := httptest.NewRecorder()

	if err := h.Upload(e.NewContext(req, rec)); err != nil {
		t.Fatalf("Upload: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestUploadSuccess(t *testing.T) {
	h, _ := newTestHandler(t)
	e := echo.New()

	body, contentType := multipartTrack(t, "song.mp3", []byte("audio"))
	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set(echo.HeaderContentType, contentType)
	rec := httptest.NewRecorder()

	if err := h.Upload(e.NewContext(req, rec)); err != nil {
		t.Fatalf("Upload: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.ID == 0 {
		t.Error("response id = 0, want non-zero")
	}
}

func TestSearch(t *testing.T) {
	h, _ := newTestHandler(t)
	e := echo.New()

	body, contentType := multipartTrack(t, "findme.mp3", []byte("audio"))
	upReq := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	upReq.Header.Set(echo.HeaderContentType, contentType)
	upRec := httptest.NewRecorder()
	if err := h.Upload(e.NewContext(upReq, upRec)); err != nil {
		t.Fatalf("Upload: %v", err)
	}
	if upRec.Code != http.StatusOK {
		t.Fatalf("upload status = %d, body=%s", upRec.Code, upRec.Body.String())
	}

	req := httptest.NewRequest(http.MethodGet, "/api/search?q=findme", nil)
	rec := httptest.NewRecorder()
	if err := h.Search(e.NewContext(req, rec)); err != nil {
		t.Fatalf("Search: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var tracks []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &tracks); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(tracks) != 1 {
		t.Fatalf("got %d tracks, want 1", len(tracks))
	}
}

func streamDirect(t *testing.T, h *TrackHandler, trackID string) *httptest.ResponseRecorder {
	t.Helper()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/stream-direct/"+trackID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/stream-direct/:track_id")
	c.SetParamNames("track_id")
	c.SetParamValues(trackID)

	if err := h.StreamDirect(c); err != nil {
		if he, ok := err.(*echo.HTTPError); ok {
			rec.Code = he.Code
		} else {
			t.Fatalf("StreamDirect: %v", err)
		}
	}
	return rec
}

func TestStreamDirectPathTraversal(t *testing.T) {
	h, _ := newTestHandler(t)

	for _, name := range []string{"..", ".", "../secret.mp3", `..\secret.mp3`, "a/b.mp3"} {
		rec := streamDirect(t, h, name)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("trackID=%q status = %d, want %d", name, rec.Code, http.StatusBadRequest)
		}
	}
}

func TestStreamDirectNotFound(t *testing.T) {
	h, _ := newTestHandler(t)

	rec := streamDirect(t, h, "missing.mp3")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestStreamDirectServesFile(t *testing.T) {
	h, uploadsDir := newTestHandler(t)

	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	filePath := filepath.Join(uploadsDir, "1.mp3")
	if err := os.WriteFile(filePath, []byte("audio-bytes"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	rec := streamDirect(t, h, "1.mp3")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	data, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if string(data) != "audio-bytes" {
		t.Errorf("body = %q, want %q", data, "audio-bytes")
	}
}

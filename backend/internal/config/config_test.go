package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("DB_PATH", "")
	t.Setenv("UPLOADS_DIR", "")
	t.Setenv("BODY_LIMIT", "")

	cfg := Load()
	if cfg.Port != ":8080" {
		t.Errorf("Port = %q, want %q", cfg.Port, ":8080")
	}
	if cfg.DBPath != "./music.db" {
		t.Errorf("DBPath = %q, want %q", cfg.DBPath, "./music.db")
	}
	if cfg.UploadsDir != "uploads" {
		t.Errorf("UploadsDir = %q, want %q", cfg.UploadsDir, "uploads")
	}
	if cfg.BodyLimit != "50M" {
		t.Errorf("BodyLimit = %q, want %q", cfg.BodyLimit, "50M")
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("PORT", ":9090")
	t.Setenv("DB_PATH", "/tmp/other.db")
	t.Setenv("UPLOADS_DIR", "/tmp/up")
	t.Setenv("BODY_LIMIT", "10M")

	cfg := Load()
	if cfg.Port != ":9090" {
		t.Errorf("Port = %q, want %q", cfg.Port, ":9090")
	}
	if cfg.DBPath != "/tmp/other.db" {
		t.Errorf("DBPath = %q, want %q", cfg.DBPath, "/tmp/other.db")
	}
	if cfg.UploadsDir != "/tmp/up" {
		t.Errorf("UploadsDir = %q, want %q", cfg.UploadsDir, "/tmp/up")
	}
	if cfg.BodyLimit != "10M" {
		t.Errorf("BodyLimit = %q, want %q", cfg.BodyLimit, "10M")
	}
}

package config

import "os"

type Config struct {
	Port       string
	DBPath     string
	UploadsDir string
	BodyLimit  string
}

func Load() Config {
	return Config{
		Port:       getEnv("PORT", ":8080"),
		DBPath:     getEnv("DB_PATH", "./music.db"),
		UploadsDir: getEnv("UPLOADS_DIR", "uploads"),
		BodyLimit:  getEnv("BODY_LIMIT", "50M"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

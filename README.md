# LocalTunes — Music Streaming Backend

Go-based REST API for uploading, searching, and streaming audio files with SQLite storage.

## Quick Start

### Prerequisites
- Go 1.22+
- SQLite (bundled)

### Run Locally

```bash
cd backend
go run ./cmd/server
```

Server starts on `http://localhost:8080`

### Docker

```bash
docker build -t localtunes .
docker run -p 8080:8080 -v $(pwd)/uploads:/app/uploads localtunes
```

## API Endpoints

### Upload Track
```bash
POST /api/upload
Content-Type: multipart/form-data

curl -F "track=@song.mp3" http://localhost:8080/api/upload
```

Supported formats:
- `.mp3`
- `.wav`
- `.ogg`
- `.flac`
- `.aac`
- `.m4a`

Response:
```json
{
  "id": 1,
  "title": "Track Title",
  "artist": "Artist Name",
  "album": "Album Name",
  "streamUrl": "/api/stream/1714571142000000000.mp3"
}
```

### Search Tracks
```bash
GET /api/search?q=artist

curl "http://localhost:8080/api/search?q=the%20beatles"
```

Response: Array of matching tracks

### Stream Audio
```bash
GET /api/stream/<filename>

curl "http://localhost:8080/api/stream/1714571142000000000.mp3" -o song.mp3
```

### Direct Range Streaming
```bash
GET /stream-direct/<filename>

curl -H "Range: bytes=0-102399" "http://localhost:8080/stream-direct/1714571142000000000.mp3"
```

This endpoint returns `206 Partial Content` for HTTP Range requests and is suitable for media players.

## Project Structure

```
backend/
├── cmd/server/main.go        # Entry point
├── internal/
│   ├── handler/              # HTTP endpoints
│   ├── service/              # Business logic
│   ├── repository/           # Database layer
│   └── models/               # Data structures
├── Dockerfile
└── go.mod
```

## Database

SQLite database with single `tracks` table:

```sql
CREATE TABLE tracks (
  id        INTEGER PRIMARY KEY,
  title     TEXT,
  artist    TEXT,
  album     TEXT,
  cover_url TEXT,
  file_name TEXT UNIQUE,
  stream_url TEXT
)
```

Database file: `music.db` (auto-created)

## Build

```bash
go build -o localtunes ./cmd/server
```

## Features

- ✅ Upload audio files (only supported audio formats)
- ✅ Search by title/artist/album
- ✅ Stream files with HTTP range requests
- ✅ Direct range streaming endpoint `/stream-direct/:track_id`
- ✅ CORS enabled
- ✅ Request logging

## Future

- User authentication
- Playlist management
- Metadata extraction
- Pagination
- Swagger API docs

## License

MIT

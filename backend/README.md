# Opinion Monitor Backend

Video opinion monitoring system backend built with Go, Gin, GORM, and MySQL.

## Prerequisites

- Go 1.21+
- MySQL 8.0+
- FFmpeg (for video processing)

## Installation

1. Install dependencies:
```bash
cd backend
go mod download
```

2. Install FFmpeg (if not already installed):
```bash
# macOS
brew install ffmpeg

# Ubuntu/Debian
sudo apt-get install ffmpeg

# Windows
# Download from https://ffmpeg.org/download.html
```

3. Create MySQL database:
```sql
CREATE DATABASE opinion_monitor CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

4. Configure application:
Copy `config.yaml` and update the settings:
- Database credentials
- OpenAI API key
- JWT secret

## Running

```bash
cd backend
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login
- `POST /api/auth/logout` - Logout
- `GET /api/auth/me` - Get current user

### Videos (Protected)
- `POST /api/videos/upload` - Upload videos (batch supported)
- `GET /api/videos` - List videos
- `GET /api/videos/:id` - Get video details
- `DELETE /api/videos/:id` - Delete video

### Reports (Protected)
- `GET /api/reports/:video_id` - Get report by video ID
- `GET /api/reports` - List all reports

### Jobs (Protected)
- `GET /api/jobs/:id/status` - Get job status
- `GET /api/jobs` - List all jobs

## Project Structure

```
backend/
├── cmd/server/          # Application entry point
├── internal/
│   ├── api/            # HTTP handlers
│   ├── models/         # Database models
│   ├── config/         # Configuration
│   └── worker/         # Job queue and workers
├── pkg/
│   ├── ai/             # OpenAI client
│   ├── auth/           # JWT utilities
│   └── video/          # Video processing
└── config.yaml         # Configuration file
```


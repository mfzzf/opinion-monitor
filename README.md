# Opinion Monitor - Video Sentiment Analysis System

A full-stack video opinion monitoring system that uses AI to extract text from video covers and perform sentiment analysis.

## Features

- 🎥 **Video Upload**: Batch video upload with drag-and-drop support
- 🤖 **AI Analysis**: Automatic text extraction from video covers using OpenAI Vision API
- 🎙️ **Audio Transcription**: Whisper large-v3 integration for audio-to-text transcription
- 📊 **Sentiment Analysis**: Detailed sentiment reports with scores, risk levels, and recommendations
- 🔐 **User Authentication**: Secure JWT-based authentication
- ⚡ **Async Processing**: Background job queue for efficient video processing
- 📱 **Modern UI**: Beautiful, responsive interface built with Next.js and shadcn/ui

## Architecture

### Backend
- **Go 1.21+** with Gin web framework
- **GORM** for database ORM
- **MySQL** for data persistence
- **FFmpeg** for video processing
- **OpenAI API** for AI-powered analysis
- **Python 3.12** with Whisper large-v3 for audio transcription

### Frontend
- **Next.js 14** with App Router
- **shadcn/ui** component library
- **TailwindCSS** for styling
- **TypeScript** for type safety

## Project Structure

```
opinion-monitor/
├── backend/                 # Go backend
│   ├── cmd/server/         # Application entry point
│   ├── internal/
│   │   ├── api/           # HTTP handlers
│   │   ├── models/        # Database models
│   │   ├── config/        # Configuration
│   │   └── worker/        # Job queue and workers
│   ├── pkg/
│   │   ├── ai/            # OpenAI client
│   │   ├── auth/          # JWT utilities
│   │   ├── video/         # Video processing
│   │   └── whisper/       # Whisper client
│   ├── whisper-service/   # Python Whisper microservice
│   ├── whisper-large-v3/  # Whisper model files
│   └── config.yaml        # Configuration file
├── frontend/               # Next.js frontend
│   ├── app/               # App router pages
│   ├── components/        # React components
│   └── lib/               # Utilities and API client
├── WHISPER_SETUP.md       # Whisper integration guide
├── QUICKSTART_WHISPER.md  # Quick start for Whisper
└── README.md
```

## Prerequisites

### Backend
- Go 1.21 or higher
- MySQL 8.0 or higher
- FFmpeg (for video processing)
- OpenAI API key (or compatible API)
- Python 3.12 with conda (for Whisper audio transcription)

### Frontend
- Node.js 18 or higher
- npm or yarn

## Installation & Setup

### 1. Database Setup

Create a MySQL database:

```sql
CREATE DATABASE opinion_monitor CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2. Backend Setup

```bash
cd backend

# Install dependencies
go mod download

# Install FFmpeg (if not already installed)
# macOS:
brew install ffmpeg

# Ubuntu/Debian:
sudo apt-get install ffmpeg

# Configure the application
# Edit config.yaml with your settings:
# - Database credentials
# - OpenAI API key
# - JWT secret

# Run the server
go run cmd/server/main.go
```

The backend will start on `http://localhost:8080`

### 3. Whisper Service Setup (Audio Transcription)

```bash
cd backend/whisper-service

# Create conda environment
conda create -n whisper-service python=3.12 -y
conda activate whisper-service

# Install dependencies
pip install -r requirements.txt

# Start service
python app.py
```

The Whisper service will start on `http://localhost:5000`

**For detailed Whisper setup instructions, see [WHISPER_SETUP.md](WHISPER_SETUP.md)**

### 4. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Configure environment
# Create .env.local file (already created):
# NEXT_PUBLIC_API_URL=http://localhost:8080

# Run the development server
npm run dev
```

The frontend will start on `http://localhost:3000`

## Configuration

### Backend Configuration (backend/config.yaml)

```yaml
server:
  port: "8080"
  upload_path: "./uploads"
  max_file_size: 524288000  # 500MB

database:
  host: "localhost"
  port: "3306"
  name: "opinion_monitor"
  user: "root"
  password: "your-password"

openai:
  api_base: "https://api.openai.com/v1"
  api_key: "sk-your-api-key"
  model_vision: "gpt-4o"
  model_chat: "gpt-4o"

jwt:
  secret: "your-secret-key-change-in-production"
  expiry: "24h"

worker:
  concurrency: 5

whisper:
  service_url: "http://localhost:5000"
```

### Frontend Configuration (frontend/.env.local)

```
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Usage

1. **Register/Login**: Create an account or login at `http://localhost:3000/login`

2. **Upload Videos**: 
   - Navigate to the Upload page
   - Drag and drop videos or click to select
   - Upload single or multiple videos at once
   - Supported formats: MP4, AVI, MOV, MKV, FLV, WMV, WebM, M4V

3. **Monitor Progress**:
   - View all videos on the Videos page
   - See real-time status updates (Pending → Processing → Completed)
   - Filter by status

4. **View Reports**:
   - Click "View Report" on completed videos
   - See extracted text, sentiment analysis, risk assessment
   - Get actionable recommendations

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login
- `POST /api/auth/logout` - Logout
- `GET /api/auth/me` - Get current user

### Videos (Protected)
- `POST /api/videos/upload` - Upload videos (batch)
- `GET /api/videos` - List videos (with pagination and filters)
- `GET /api/videos/:id` - Get video details
- `DELETE /api/videos/:id` - Delete video

### Reports (Protected)
- `GET /api/reports/:video_id` - Get report by video ID
- `GET /api/reports` - List all reports (with pagination)

### Jobs (Protected)
- `GET /api/jobs/:id/status` - Get job status
- `GET /api/jobs` - List jobs (with pagination)

## How It Works

1. **Upload**: User uploads video(s) through the web interface
2. **Storage**: Videos are saved to local disk organized by user and date
3. **Job Creation**: A job is created for each video and queued
4. **Processing**: Worker pool picks up jobs asynchronously:
   - Extract cover frame from video at 1 second using FFmpeg
   - Extract audio from video and transcribe using Whisper large-v3
   - Use OpenAI Vision API to extract text from cover image
   - Combine cover text and audio transcription
   - Use OpenAI Chat API to analyze sentiment on combined text
   - Save results to database with both text sources
5. **Display**: User views detailed sentiment analysis report with cover text and transcript

## Development

### Backend Development

```bash
cd backend
go run cmd/server/main.go

# Run with hot reload (install air first)
air
```

### Frontend Development

```bash
cd frontend
npm run dev

# Build for production
npm run build
npm run start
```

## Production Deployment

### Backend

```bash
cd backend
go build -o opinion-monitor cmd/server/main.go
./opinion-monitor
```

### Frontend

```bash
cd frontend
npm run build
npm run start
```

## Security Considerations

- Change JWT secret in production
- Use environment variables for sensitive data
- Enable HTTPS in production
- Set up proper CORS policies
- Implement rate limiting
- Regular security updates

## Troubleshooting

### FFmpeg not found
Make sure FFmpeg is installed and available in PATH:
```bash
ffmpeg -version
```

### Database connection failed
Check MySQL is running and credentials are correct in config.yaml

### OpenAI API errors
Verify your API key is valid and has sufficient credits

### CORS issues
Ensure frontend URL is in the CORS allowed origins in backend/cmd/server/main.go

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.


# Quick Start Guide

Get the Opinion Monitor system up and running in 5 minutes.

## Prerequisites Check

Before starting, ensure you have:

```bash
# Check Go version (need 1.21+)
go version

# Check Node.js version (need 18+)
node --version

# Check MySQL is running
mysql --version

# Check FFmpeg is installed
ffmpeg -version
```

If anything is missing:

```bash
# Install FFmpeg
# macOS:
brew install ffmpeg

# Ubuntu/Debian:
sudo apt-get install ffmpeg

# Windows: Download from https://ffmpeg.org/download.html
```

## Step 1: Database Setup (2 minutes)

```bash
# Login to MySQL
mysql -u root -p

# Create database
CREATE DATABASE opinion_monitor CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# Exit MySQL
exit
```

## Step 2: Backend Setup (2 minutes)

```bash
cd backend

# Download dependencies
go mod download

# Edit config.yaml with your settings
# IMPORTANT: Update these fields:
# - database.password (your MySQL password)
# - openai.api_key (your OpenAI API key)
# - jwt.secret (change to a random string)

# Start the backend
go run cmd/server/main.go
```

You should see: `Server starting on port 8080`

## Step 3: Frontend Setup (1 minute)

Open a new terminal:

```bash
cd frontend

# Install dependencies
npm install

# Start the frontend
npm run dev
```

You should see: `Ready on http://localhost:3000`

## Step 4: Use the System

1. Open browser to http://localhost:3000
2. Click "Register" and create an account
3. Upload a video from the Upload page
4. Watch it process (refresh Videos page to see status)
5. Click "View Report" when completed

## Troubleshooting

### Backend won't start
- Check MySQL is running: `mysql -u root -p`
- Verify database exists: `SHOW DATABASES;`
- Check config.yaml has correct credentials

### Frontend won't start
- Delete node_modules and try again: `rm -rf node_modules && npm install`
- Check port 3000 is not in use

### Video processing fails
- Check FFmpeg is installed: `ffmpeg -version`
- Verify OpenAI API key is valid
- Check backend logs for errors

### Upload fails
- Check video file size (default limit: 500MB)
- Ensure video is in supported format (MP4, AVI, MOV, etc.)
- Check backend is running on port 8080

## Configuration Tips

### For Development
Current settings in config.yaml are fine.

### For Production
Update these in config.yaml:

```yaml
jwt:
  secret: "use-a-long-random-string-here"

openai:
  api_key: "your-production-api-key"

server:
  upload_path: "/var/opinion-monitor/uploads"
```

## What's Next?

- Read the full [README.md](README.md) for detailed documentation
- Check the API documentation in [backend/README.md](backend/README.md)
- Explore frontend customization in [frontend/README.md](frontend/README.md)

## Need Help?

Common issues and solutions:

1. **"connection refused"** → Backend not running
2. **"401 Unauthorized"** → Login again
3. **"failed to extract cover"** → FFmpeg not installed
4. **"API request failed"** → Check OpenAI API key and credits

## Test Data

For testing, use any short video file (< 1 minute recommended for faster processing).

Good test videos:
- Screen recordings
- Short clips from social media
- Sample videos from stock footage sites

## Performance

- First video upload: May take 10-30 seconds to process
- Subsequent uploads: Processed in parallel (up to 5 concurrent jobs)
- Frontend auto-refreshes status every 5 seconds

Enjoy using Opinion Monitor! 🎥📊


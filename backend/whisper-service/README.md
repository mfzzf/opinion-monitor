# Whisper Audio Transcription Microservice

Python microservice for transcribing video audio using Whisper large-v3 model.

## Setup

### Using Conda (Python 3.12)

```bash
# Create conda environment
conda create -n whisper-service python=3.12 -y
conda activate whisper-service

# Install dependencies
pip install -r requirements.txt
```

### Running the Service

```bash
conda activate whisper-service
python app.py
```

The service will start on `http://localhost:5000`

## API Endpoints

### Health Check

```bash
GET /health
```

Response:
```json
{
  "status": "healthy",
  "model_loaded": true,
  "device": "cpu"
}
```

### Transcribe Video

```bash
POST /transcribe
Content-Type: application/json

{
  "video_path": "/path/to/video.mp4"
}
```

Response:
```json
{
  "success": true,
  "transcription": "transcribed text from video audio",
  "language": "zh"
}
```

## Requirements

- Python 3.12
- ffmpeg (must be installed on system)
- CUDA-compatible GPU (optional, for faster inference)

## Model

The service uses Whisper large-v3 model located at `../whisper-large-v3/`


#!/usr/bin/env python3
"""
Whisper Audio Transcription Microservice
Provides HTTP API for transcribing audio from video files using Whisper large-v3
"""

from flask import Flask, request, jsonify
import torch
from transformers import AutoModelForSpeechSeq2Seq, AutoProcessor, pipeline
import os
import subprocess
import tempfile
import logging
from pathlib import Path

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = Flask(__name__)

# Global variables for model
whisper_pipeline = None
MODEL_PATH = "../whisper-large-v3"

def initialize_model():
    """Initialize Whisper model on startup"""
    global whisper_pipeline
    
    logger.info("Initializing Whisper large-v3 model...")
    
    device = "cuda:0" if torch.cuda.is_available() else "cpu"
    torch_dtype = torch.float16 if torch.cuda.is_available() else torch.float32
    
    logger.info(f"Using device: {device}, dtype: {torch_dtype}")
    
    try:
        # Load model from local directory
        model = AutoModelForSpeechSeq2Seq.from_pretrained(
            MODEL_PATH, 
            torch_dtype=torch_dtype, 
            low_cpu_mem_usage=True, 
            use_safetensors=True
        )
        model.to(device)
        
        processor = AutoProcessor.from_pretrained(MODEL_PATH)
        
        whisper_pipeline = pipeline(
            "automatic-speech-recognition",
            model=model,
            tokenizer=processor.tokenizer,
            feature_extractor=processor.feature_extractor,
            torch_dtype=torch_dtype,
            device=device,
        )
        
        logger.info("Whisper model initialized successfully")
    except Exception as e:
        logger.error(f"Failed to initialize model: {e}")
        raise

def extract_audio_from_video(video_path, output_audio_path):
    """Extract audio from video file using ffmpeg"""
    try:
        cmd = [
            "ffmpeg",
            "-i", video_path,
            "-vn",  # No video
            "-acodec", "pcm_s16le",  # PCM 16-bit little-endian
            "-ar", "16000",  # 16kHz sample rate
            "-ac", "1",  # Mono
            "-y",  # Overwrite output file
            output_audio_path
        ]
        
        result = subprocess.run(
            cmd, 
            capture_output=True, 
            text=True,
            check=True
        )
        
        logger.info(f"Audio extracted successfully to {output_audio_path}")
        return True
    except subprocess.CalledProcessError as e:
        logger.error(f"ffmpeg failed: {e.stderr}")
        return False
    except Exception as e:
        logger.error(f"Error extracting audio: {e}")
        return False

@app.route('/health', methods=['GET'])
def health_check():
    """Health check endpoint"""
    return jsonify({
        "status": "healthy",
        "model_loaded": whisper_pipeline is not None,
        "device": "cuda" if torch.cuda.is_available() else "cpu"
    })

@app.route('/transcribe', methods=['POST'])
def transcribe():
    """
    Transcribe audio from video file
    
    Request JSON:
    {
        "video_path": "/path/to/video.mp4"
    }
    
    Response JSON:
    {
        "success": true,
        "transcription": "transcribed text...",
        "language": "zh"
    }
    """
    if whisper_pipeline is None:
        return jsonify({
            "success": False,
            "error": "Model not initialized"
        }), 503
    
    try:
        # Get video path from request
        data = request.get_json()
        video_path = data.get('video_path')
        
        if not video_path:
            return jsonify({
                "success": False,
                "error": "video_path is required"
            }), 400
        
        if not os.path.exists(video_path):
            return jsonify({
                "success": False,
                "error": f"Video file not found: {video_path}"
            }), 404
        
        logger.info(f"Processing video: {video_path}")
        
        # Create temporary file for audio
        with tempfile.NamedTemporaryFile(suffix='.wav', delete=False) as temp_audio:
            temp_audio_path = temp_audio.name
        
        try:
            # Extract audio from video
            if not extract_audio_from_video(video_path, temp_audio_path):
                return jsonify({
                    "success": False,
                    "error": "Failed to extract audio from video"
                }), 500
            
            # Check if audio file has content
            if os.path.getsize(temp_audio_path) < 1000:
                logger.warning("Audio file is too small, video might have no audio")
                return jsonify({
                    "success": True,
                    "transcription": "",
                    "language": "unknown",
                    "warning": "No audio detected in video"
                })
            
            # Transcribe audio
            logger.info("Starting transcription...")
            result = whisper_pipeline(
                temp_audio_path,
                return_timestamps=True,  # Required for videos > 30 seconds
                generate_kwargs={
                    "language": "chinese",  # Set to Chinese by default, can be made configurable
                    "task": "transcribe"
                }
            )
            
            transcription = result["text"].strip()
            logger.info(f"Transcription completed: {transcription[:100]}...")
            
            return jsonify({
                "success": True,
                "transcription": transcription,
                "language": "zh"
            })
            
        finally:
            # Clean up temporary audio file
            if os.path.exists(temp_audio_path):
                os.remove(temp_audio_path)
                logger.info(f"Cleaned up temporary file: {temp_audio_path}")
    
    except Exception as e:
        logger.error(f"Transcription error: {e}", exc_info=True)
        return jsonify({
            "success": False,
            "error": str(e)
        }), 500

if __name__ == '__main__':
    # Initialize model on startup
    initialize_model()
    
    # Start Flask server
    app.run(host='0.0.0.0', port=5000, debug=False)


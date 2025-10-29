#!/usr/bin/env python3
"""
Test script for Whisper service
"""

import requests
import json
import sys

def test_health():
    """Test health endpoint"""
    print("Testing /health endpoint...")
    try:
        response = requests.get("http://localhost:5000/health", timeout=5)
        response.raise_for_status()
        data = response.json()
        print(f"✓ Health check passed: {json.dumps(data, indent=2)}")
        return True
    except Exception as e:
        print(f"✗ Health check failed: {e}")
        return False

def test_transcribe(video_path):
    """Test transcribe endpoint"""
    print(f"\nTesting /transcribe endpoint with: {video_path}")
    try:
        payload = {"video_path": video_path}
        response = requests.post(
            "http://localhost:5000/transcribe",
            json=payload,
            timeout=300  # 5 minutes
        )
        response.raise_for_status()
        data = response.json()
        
        if data.get("success"):
            print(f"✓ Transcription successful!")
            print(f"  Language: {data.get('language', 'unknown')}")
            print(f"  Transcription preview: {data.get('transcription', '')[:200]}...")
            return True
        else:
            print(f"✗ Transcription failed: {data.get('error', 'unknown error')}")
            return False
    except Exception as e:
        print(f"✗ Transcription request failed: {e}")
        return False

if __name__ == "__main__":
    print("=" * 60)
    print("Whisper Service Test Script")
    print("=" * 60)
    
    # Test health
    if not test_health():
        print("\n⚠ Service is not healthy. Make sure it's running:")
        print("  python app.py")
        sys.exit(1)
    
    # Test transcription (if video path provided)
    if len(sys.argv) > 1:
        video_path = sys.argv[1]
        test_transcribe(video_path)
    else:
        print("\n✓ All tests passed!")
        print("\nTo test transcription, provide a video path:")
        print(f"  python {sys.argv[0]} /path/to/video.mp4")


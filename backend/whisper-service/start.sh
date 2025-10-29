#!/bin/bash
# Start script for Whisper microservice

echo "Starting Whisper Audio Transcription Service..."

# Check if conda is available
if ! command -v conda &> /dev/null; then
    echo "Error: conda is not installed or not in PATH"
    exit 1
fi

# Check if environment exists
if ! conda env list | grep -q "whisper-service"; then
    echo "Creating conda environment..."
    conda create -n whisper-service python=3.12 -y
    
    echo "Installing dependencies..."
    conda activate whisper-service
    pip install -r requirements.txt
else
    echo "Conda environment already exists"
fi

# Activate environment and start service
echo "Activating conda environment..."
eval "$(conda shell.bash hook)"
conda activate whisper-service

echo "Starting Flask service on port 5000..."
python app.py


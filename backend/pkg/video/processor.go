package video

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Processor struct{}

func NewProcessor() *Processor {
	return &Processor{}
}

// ExtractCover extracts a frame from the video at the specified timestamp (in seconds)
func (p *Processor) ExtractCover(videoPath, outputPath string, timestamp float64) error {
	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Use ffmpeg to extract frame
	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-ss", fmt.Sprintf("%.2f", timestamp),
		"-vframes", "1",
		"-q:v", "2",
		"-y",
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w, output: %s", err, string(output))
	}

	return nil
}

// GetVideoDuration gets the duration of a video in seconds
func (p *Processor) GetVideoDuration(videoPath string) (float64, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		videoPath,
	)

	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe failed: %w", err)
	}

	var duration float64
	_, err = fmt.Sscanf(strings.TrimSpace(string(output)), "%f", &duration)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %w", err)
	}

	return duration, nil
}

// ExtractAudio extracts audio from the video file
func (p *Processor) ExtractAudio(videoPath, outputPath string) error {
	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Use ffmpeg to extract audio to WAV format
	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-vn",                    // No video
		"-acodec", "pcm_s16le",   // PCM 16-bit little-endian
		"-ar", "16000",           // 16kHz sample rate
		"-ac", "1",               // Mono
		"-y",                     // Overwrite output file
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w, output: %s", err, string(output))
	}

	return nil
}

// IsVideoFile checks if the file is a supported video format
func IsVideoFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	supportedExts := []string{".mp4", ".avi", ".mov", ".mkv", ".flv", ".wmv", ".webm", ".m4v"}

	for _, supportedExt := range supportedExts {
		if ext == supportedExt {
			return true
		}
	}

	return false
}

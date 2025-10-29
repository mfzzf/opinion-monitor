package whisper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type TranscribeRequest struct {
	VideoPath string `json:"video_path"`
}

type TranscribeResponse struct {
	Success       bool   `json:"success"`
	Transcription string `json:"transcription"`
	Language      string `json:"language"`
	Error         string `json:"error,omitempty"`
	Warning       string `json:"warning,omitempty"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Minute, // Long timeout for transcription
		},
	}
}

// TranscribeAudio sends a video file path to the Whisper service for transcription
func (c *Client) TranscribeAudio(videoPath string) (string, error) {
	// Prepare request
	reqBody := TranscribeRequest{
		VideoPath: videoPath,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make HTTP request
	url := fmt.Sprintf("%s/transcribe", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Whisper service: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var transcribeResp TranscribeResponse
	if err := json.Unmarshal(body, &transcribeResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w, body: %s", err, string(body))
	}

	if !transcribeResp.Success {
		return "", fmt.Errorf("transcription failed: %s", transcribeResp.Error)
	}

	return transcribeResp.Transcription, nil
}

// HealthCheck checks if the Whisper service is available
func (c *Client) HealthCheck() error {
	url := fmt.Sprintf("%s/health", c.baseURL)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check returned status %d", resp.StatusCode)
	}

	return nil
}


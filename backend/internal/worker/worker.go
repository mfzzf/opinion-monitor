package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"opinion-monitor/internal/config"
	"opinion-monitor/internal/models"
	"opinion-monitor/pkg/ai"
	"opinion-monitor/pkg/video"
	"opinion-monitor/pkg/whisper"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

type WorkerPool struct {
	cfg           *config.Config
	db            *gorm.DB
	queue         *JobQueue
	aiClient      *ai.OpenAIClient
	whisperClient *whisper.Client
	processor     *video.Processor
}

func NewWorkerPool(cfg *config.Config, db *gorm.DB, queue *JobQueue, whisperClient *whisper.Client) *WorkerPool {
	aiClient := ai.NewOpenAIClient(
		cfg.OpenAI.APIBase,
		cfg.OpenAI.APIKey,
		cfg.OpenAI.ModelVision,
		cfg.OpenAI.ModelChat,
	)

	return &WorkerPool{
		cfg:           cfg,
		db:            db,
		queue:         queue,
		aiClient:      aiClient,
		whisperClient: whisperClient,
		processor:     video.NewProcessor(),
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.cfg.Worker.Concurrency; i++ {
		go wp.worker(i)
	}
	log.Printf("Started %d workers", wp.cfg.Worker.Concurrency)
}

func (wp *WorkerPool) worker(id int) {
	log.Printf("Worker %d started", id)

	for videoID := range wp.queue.Jobs() {
		log.Printf("Worker %d processing video %d", id, videoID)

		if err := wp.processVideo(videoID); err != nil {
			log.Printf("Worker %d failed to process video %d: %v", id, videoID, err)
			wp.markJobFailed(videoID, err.Error())
		}
	}
}

func (wp *WorkerPool) processVideo(videoID uint) error {
	startTime := time.Now()

	// Update job status to processing
	if err := wp.updateJobStatus(videoID, models.JobStatusProcessing, ""); err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	// Update video status to processing
	if err := wp.updateVideoStatus(videoID, models.StatusProcessing); err != nil {
		return fmt.Errorf("failed to update video status: %w", err)
	}

	// Get video info
	var videoRecord models.Video
	if err := wp.db.First(&videoRecord, videoID).Error; err != nil {
		return fmt.Errorf("failed to get video: %w", err)
	}

	// Extract cover frame
	coverFilename := fmt.Sprintf("cover_%d.jpg", videoID)
	coverPath := filepath.Join(filepath.Dir(videoRecord.FilePath), coverFilename)

	if err := wp.processor.ExtractCover(videoRecord.FilePath, coverPath, 1.0); err != nil {
		return fmt.Errorf("failed to extract cover: %w", err)
	}

	// Update video with cover path
	if err := wp.db.Model(&videoRecord).Update("cover_path", coverPath).Error; err != nil {
		return fmt.Errorf("failed to update cover path: %w", err)
	}

	// Extract audio from video
	audioFilename := fmt.Sprintf("audio_%d.wav", videoID)
	audioPath := filepath.Join(filepath.Dir(videoRecord.FilePath), audioFilename)

	if err := wp.processor.ExtractAudio(videoRecord.FilePath, audioPath); err != nil {
		log.Printf("Warning: failed to extract audio: %v", err)
		// Continue processing even if audio extraction fails
	}

	// Update video with audio path
	if err := wp.db.Model(&videoRecord).Update("audio_path", audioPath).Error; err != nil {
		return fmt.Errorf("failed to update audio path: %w", err)
	}

	// Transcribe audio using Whisper
	var transcriptText string
	if wp.whisperClient != nil {
		// Get absolute path for video file
		absVideoPath := videoRecord.FilePath
		if !filepath.IsAbs(absVideoPath) {
			// If relative path, prepend current working directory
			cwd, err := os.Getwd()
			if err == nil {
				absVideoPath = filepath.Join(cwd, absVideoPath)
			}
		}
		
		transcript, err := wp.whisperClient.TranscribeAudio(absVideoPath)
		if err != nil {
			log.Printf("Warning: failed to transcribe audio: %v", err)
			// Continue processing even if transcription fails
		} else {
			transcriptText = transcript
			log.Printf("Transcribed text: %s", transcriptText)

			// Update video with transcript text
			if err := wp.db.Model(&videoRecord).Update("transcript_text", transcriptText).Error; err != nil {
				return fmt.Errorf("failed to update transcript text: %w", err)
			}
		}
	}

	// Extract text from cover using AI
	coverText, err := wp.aiClient.ExtractTextFromImage(coverPath)
	if err != nil {
		return fmt.Errorf("failed to extract text from image: %w", err)
	}

	log.Printf("Extracted cover text: %s", coverText)

	// Combine cover text and transcript for comprehensive analysis
	var combinedText strings.Builder
	combinedText.WriteString("封面文字：\n")
	combinedText.WriteString(coverText)
	if transcriptText != "" {
		combinedText.WriteString("\n\n音频转录文字：\n")
		combinedText.WriteString(transcriptText)
	}

	// Analyze sentiment using combined text
	report, err := wp.aiClient.AnalyzeSentiment(combinedText.String())
	if err != nil {
		return fmt.Errorf("failed to analyze sentiment: %w", err)
	}

	processingTime := time.Since(startTime).Seconds()

	// Convert arrays to JSON strings
	keyTopicsJSON, _ := json.Marshal(report.KeyTopics)
	recommendationsJSON, _ := json.Marshal(report.Recommendations)

	// Save report
	reportRecord := models.Report{
		VideoID:          videoID,
		CoverText:        coverText,
		TranscriptText:   transcriptText,
		SentimentScore:   report.SentimentScore,
		SentimentLabel:   report.SentimentLabel,
		KeyTopics:        string(keyTopicsJSON),
		RiskLevel:        report.RiskLevel,
		DetailedAnalysis: report.DetailedAnalysis,
		Recommendations:  string(recommendationsJSON),
		ProcessingTime:   processingTime,
	}

	if err := wp.db.Create(&reportRecord).Error; err != nil {
		return fmt.Errorf("failed to save report: %w", err)
	}

	// Update video status to completed
	if err := wp.updateVideoStatus(videoID, models.StatusCompleted); err != nil {
		return fmt.Errorf("failed to update video status: %w", err)
	}

	// Update job status to completed
	if err := wp.updateJobStatus(videoID, models.JobStatusCompleted, ""); err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	log.Printf("Successfully processed video %d in %.2f seconds", videoID, processingTime)
	return nil
}

func (wp *WorkerPool) updateVideoStatus(videoID uint, status models.VideoStatus) error {
	return wp.db.Model(&models.Video{}).Where("id = ?", videoID).Update("status", status).Error
}

func (wp *WorkerPool) updateJobStatus(videoID uint, status models.JobStatus, errorMsg string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if errorMsg != "" {
		updates["error_message"] = errorMsg
	}

	return wp.db.Model(&models.Job{}).Where("video_id = ?", videoID).Updates(updates).Error
}

func (wp *WorkerPool) markJobFailed(videoID uint, errorMsg string) {
	wp.updateVideoStatus(videoID, models.StatusFailed)
	wp.updateJobStatus(videoID, models.JobStatusFailed, errorMsg)
}

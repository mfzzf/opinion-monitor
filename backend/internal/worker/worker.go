package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"opinion-monitor/internal/config"
	"opinion-monitor/internal/models"
	"opinion-monitor/pkg/ai"
	"opinion-monitor/pkg/video"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

type WorkerPool struct {
	cfg       *config.Config
	db        *gorm.DB
	queue     *JobQueue
	aiClient  *ai.OpenAIClient
	processor *video.Processor
}

func NewWorkerPool(cfg *config.Config, db *gorm.DB, queue *JobQueue) *WorkerPool {
	aiClient := ai.NewOpenAIClient(
		cfg.OpenAI.APIBase,
		cfg.OpenAI.APIKey,
		cfg.OpenAI.ModelVision,
		cfg.OpenAI.ModelChat,
	)

	return &WorkerPool{
		cfg:       cfg,
		db:        db,
		queue:     queue,
		aiClient:  aiClient,
		processor: video.NewProcessor(),
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

	// Extract text from cover using AI
	coverText, err := wp.aiClient.ExtractTextFromImage(coverPath)
	if err != nil {
		return fmt.Errorf("failed to extract text from image: %w", err)
	}

	log.Printf("Extracted text: %s", coverText)

	// Analyze sentiment
	report, err := wp.aiClient.AnalyzeSentiment(coverText)
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

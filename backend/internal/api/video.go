package api

import (
	"fmt"
	"net/http"
	"opinion-monitor/internal/config"
	"opinion-monitor/internal/models"
	"opinion-monitor/internal/worker"
	"opinion-monitor/pkg/video"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VideoHandler struct {
	db       *gorm.DB
	cfg      *config.Config
	jobQueue *worker.JobQueue
}

func NewVideoHandler(db *gorm.DB, cfg *config.Config, jobQueue *worker.JobQueue) *VideoHandler {
	return &VideoHandler{
		db:       db,
		cfg:      cfg,
		jobQueue: jobQueue,
	}
}

func (h *VideoHandler) Upload(c *gin.Context) {
	userID, _ := c.Get("user_id")

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	files := form.File["videos"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files provided"})
		return
	}

	// Create upload directory
	userDir := filepath.Join(h.cfg.Server.UploadPath, fmt.Sprintf("%d", userID), time.Now().Format("2006-01-02"))
	if err := os.MkdirAll(userDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	uploadedVideos := []models.Video{}

	for _, file := range files {
		// Validate file type
		if !video.IsVideoFile(file.Filename) {
			continue
		}

		// Check file size
		if file.Size > h.cfg.Server.MaxFileSize {
			continue
		}

		// Generate unique filename
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
		filePath := filepath.Join(userDir, filename)

		// Save file
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			continue
		}

		// Get video duration
		processor := video.NewProcessor()
		duration, _ := processor.GetVideoDuration(filePath)

		// Create video record
		videoRecord := models.Video{
			UserID:           userID.(uint),
			OriginalFilename: file.Filename,
			FilePath:         filePath,
			FileSize:         file.Size,
			Duration:         duration,
			Status:           models.StatusPending,
		}

		if err := h.db.Create(&videoRecord).Error; err != nil {
			os.Remove(filePath)
			continue
		}

		// Create job
		job := models.Job{
			VideoID: videoRecord.ID,
			Status:  models.JobStatusPending,
		}

		if err := h.db.Create(&job).Error; err != nil {
			continue
		}

		// Queue job
		h.jobQueue.Push(videoRecord.ID)

		uploadedVideos = append(uploadedVideos, videoRecord)
	}

	if len(uploadedVideos) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid videos uploaded"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Videos uploaded successfully",
		"count":   len(uploadedVideos),
		"videos":  uploadedVideos,
	})
}

func (h *VideoHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	query := h.db.Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Model(&models.Video{}).Count(&total)

	var videos []models.Video
	if err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&videos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch videos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"videos":    videos,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *VideoHandler) Get(c *gin.Context) {
	userID, _ := c.Get("user_id")
	videoID := c.Param("id")

	var videoRecord models.Video
	if err := h.db.Where("id = ? AND user_id = ?", videoID, userID).
		Preload("Report").
		Preload("Job").
		First(&videoRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch video"})
		return
	}

	c.JSON(http.StatusOK, videoRecord)
}

func (h *VideoHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	videoID := c.Param("id")

	var videoRecord models.Video
	if err := h.db.Where("id = ? AND user_id = ?", videoID, userID).First(&videoRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch video"})
		return
	}

	// Delete files
	os.Remove(videoRecord.FilePath)
	if videoRecord.CoverPath != "" {
		os.Remove(videoRecord.CoverPath)
	}

	// Delete from database (soft delete)
	if err := h.db.Delete(&videoRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete video"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video deleted successfully"})
}

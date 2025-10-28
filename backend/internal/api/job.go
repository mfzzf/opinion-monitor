package api

import (
	"net/http"
	"opinion-monitor/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type JobHandler struct {
	db *gorm.DB
}

func NewJobHandler(db *gorm.DB) *JobHandler {
	return &JobHandler{db: db}
}

func (h *JobHandler) GetStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	jobID := c.Param("id")

	// Check if job belongs to user
	var job models.Job
	if err := h.db.Preload("Video").Where("id = ?", jobID).First(&job).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch job"})
		return
	}

	if job.Video.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, job)
}

func (h *JobHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var jobs []models.Job
	var total int64

	h.db.Model(&models.Job{}).
		Joins("JOIN videos ON jobs.video_id = videos.id").
		Where("videos.user_id = ?", userID).
		Count(&total)

	if err := h.db.Preload("Video").
		Joins("JOIN videos ON jobs.video_id = videos.id").
		Where("videos.user_id = ?", userID).
		Order("jobs.created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs":      jobs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

package api

import (
	"net/http"
	"opinion-monitor/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReportHandler struct {
	db *gorm.DB
}

func NewReportHandler(db *gorm.DB) *ReportHandler {
	return &ReportHandler{db: db}
}

func (h *ReportHandler) GetByVideoID(c *gin.Context) {
	userID, _ := c.Get("user_id")
	videoID := c.Param("video_id")

	// Check if video belongs to user
	var video models.Video
	if err := h.db.Where("id = ? AND user_id = ?", videoID, userID).First(&video).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch video"})
		return
	}

	var report models.Report
	if err := h.db.Where("video_id = ?", videoID).Preload("Video").First(&report).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch report"})
		return
	}

	c.JSON(http.StatusOK, report)
}

func (h *ReportHandler) List(c *gin.Context) {
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

	// Join with videos table to filter by user_id
	var reports []models.Report
	var total int64

	h.db.Model(&models.Report{}).
		Joins("JOIN videos ON reports.video_id = videos.id").
		Where("videos.user_id = ?", userID).
		Count(&total)

	if err := h.db.Preload("Video").
		Joins("JOIN videos ON reports.video_id = videos.id").
		Where("videos.user_id = ?", userID).
		Order("reports.created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reports"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reports":   reports,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

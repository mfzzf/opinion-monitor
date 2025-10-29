package models

import (
	"time"

	"gorm.io/gorm"
)

type Report struct {
	ID               uint           `gorm:"primarykey" json:"id"`
	VideoID          uint           `gorm:"uniqueIndex;not null" json:"video_id"`
	CoverText        string         `gorm:"type:text" json:"cover_text"`
	TranscriptText   string         `gorm:"type:text" json:"transcript_text"`
	SentimentScore   float64        `json:"sentiment_score"`
	SentimentLabel   string         `gorm:"type:varchar(20)" json:"sentiment_label"`
	KeyTopics        string         `gorm:"type:text" json:"key_topics"` // JSON array stored as string
	RiskLevel        string         `gorm:"type:varchar(20)" json:"risk_level"`
	DetailedAnalysis string         `gorm:"type:text" json:"detailed_analysis"`
	Recommendations  string         `gorm:"type:text" json:"recommendations"` // JSON array stored as string
	ProcessingTime   float64        `json:"processing_time"`                  // in seconds
	CreatedAt        time.Time      `json:"created_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	Video Video `gorm:"foreignKey:VideoID" json:"video,omitempty"`
}

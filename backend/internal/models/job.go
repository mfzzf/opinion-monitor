package models

import (
	"time"

	"gorm.io/gorm"
)

type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

type Job struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	VideoID      uint           `gorm:"uniqueIndex;not null" json:"video_id"`
	Status       JobStatus      `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	RetryCount   int            `gorm:"default:0" json:"retry_count"`
	ErrorMessage string         `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Video Video `gorm:"foreignKey:VideoID" json:"video,omitempty"`
}

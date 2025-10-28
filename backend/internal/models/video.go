package models

import (
	"time"

	"gorm.io/gorm"
)

type VideoStatus string

const (
	StatusPending    VideoStatus = "pending"
	StatusProcessing VideoStatus = "processing"
	StatusCompleted  VideoStatus = "completed"
	StatusFailed     VideoStatus = "failed"
)

type Video struct {
	ID               uint           `gorm:"primarykey" json:"id"`
	UserID           uint           `gorm:"not null;index" json:"user_id"`
	OriginalFilename string         `gorm:"type:varchar(255);not null" json:"original_filename"`
	FilePath         string         `gorm:"type:varchar(500);not null" json:"file_path"`
	CoverPath        string         `gorm:"type:varchar(500)" json:"cover_path"`
	FileSize         int64          `json:"file_size"`
	Duration         float64        `json:"duration"`
	Status           VideoStatus    `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	User   User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Report *Report `gorm:"foreignKey:VideoID" json:"report,omitempty"`
	Job    *Job    `gorm:"foreignKey:VideoID" json:"job,omitempty"`
}

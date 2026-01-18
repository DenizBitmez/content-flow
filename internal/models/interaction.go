package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ContentID uint           `gorm:"index" json:"content_id"`
	UserID    uint           `gorm:"index" json:"user_id"`
	User      User           `json:"user"` // Preload user info
	Body      string         `json:"body"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Like represents a unique like per user per content
type Like struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	ContentID uint      `gorm:"primaryKey" json:"content_id"`
	CreatedAt time.Time `json:"created_at"`
}

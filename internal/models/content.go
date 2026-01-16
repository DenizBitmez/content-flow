package models

import (
	"time"

	"gorm.io/gorm"
)

type Content struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Title      string         `json:"title"`
	Slug       string         `gorm:"uniqueIndex" json:"slug"`
	Body       string         `json:"body"`
	Type       string         `json:"type"`       // e.g "Product", "Blog"
	Attributes string         `json:"attributes"` // JSON string for flexible data
	Status     string         `json:"status"`     // DRAFT, PUBLISHED
	Version    int            `json:"version"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type ContentVersion struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ContentID  uint      `gorm:"index" json:"content_id"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	Type       string    `json:"type"`
	Attributes string    `json:"attributes"`
	Status     string    `json:"status"`
	Version    int       `json:"version"`
	ChangedAt  time.Time `json:"changed_at"`
}

type ContentUpdateRequest struct {
	Title      string `json:"title"`
	Body       string `json:"body"`
	Type       string `json:"type"`
	Attributes string `json:"attributes"`
	Status     string `json:"status"`
}

type ContentCreateRequest struct {
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	Body       string `json:"body"`
	Type       string `json:"type"`
	Attributes string `json:"attributes"`
	Status     string `json:"status"`
}

package models

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Content struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `json:"title"`
	Slug        string         `gorm:"uniqueIndex:idx_slug_lang" json:"slug"`
	Body        string         `json:"body"`
	Type        string         `json:"type"`       // e.g "Product", "Blog"
	Attributes  string         `json:"attributes"` // JSON string for flexible data
	Status      string         `json:"status"`     // DRAFT, PUBLISHED
	Language    string         `gorm:"default:'en';uniqueIndex:idx_slug_lang" json:"language"`
	GroupID     string         `gorm:"index" json:"group_id"` // UUID to link translations (same content, diff lang)
	Version     int            `json:"version"`
	Categories  []Category     `gorm:"many2many:content_categories;" json:"categories,omitempty"`
	Tags        []Tag          `gorm:"many2many:content_tags;" json:"tags,omitempty"`
	AuthorID    uint           `gorm:"index" json:"author_id"`
	Author      User           `json:"author,omitempty"`
	PublishedAt *time.Time     `json:"published_at"`
	Blocks      datatypes.JSON `json:"blocks" swaggertype:"object"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type ContentVersion struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ContentID  uint           `gorm:"index" json:"content_id"`
	Title      string         `json:"title"`
	Body       string         `json:"body"`
	Blocks     datatypes.JSON `json:"blocks" swaggertype:"object"`
	Type       string         `json:"type"`
	Attributes string         `json:"attributes"`
	Status     string         `json:"status"`
	Language   string         `json:"language"`
	Version    int            `json:"version"`
	ChangedAt  time.Time      `json:"changed_at"`
}

type ContentUpdateRequest struct {
	Title string `json:"title" validate:"omitempty,min=3"` // omitempty allows partial updates if we handled PATCH, but for PUT usually full replace?
	// Actually, PUT usually means full replace. But let's assume we might want optional.
	// If it's PUT, usually all fields are expected or they get zeroed.
	// Let's stick to strict validation for PUT or check usage.
	// Since struct is used for generic update, let's just allow omitempty for flexibility or require if it's strictly PUT.
	// Given previous update logic: services.UpdateContent takes all args.
	// Let's add standard validation.
	Body        string          `json:"body"`
	Blocks      json.RawMessage `json:"blocks" swaggertype:"object"`
	Type        string          `json:"type" validate:"omitempty"`
	Attributes  string          `json:"attributes"`
	Status      string          `json:"status" validate:"omitempty,oneof=DRAFT PUBLISHED SCHEDULED"`
	Language    string          `json:"language" validate:"omitempty,len=2"`
	CategoryIDs []uint          `json:"category_ids"`
	Tags        []string        `json:"tags"` // Tag names
	PublishedAt *time.Time      `json:"published_at"`
}

type ContentCreateRequest struct {
	Title       string          `json:"title" validate:"required,min=3"`
	Slug        string          `json:"slug" validate:"required,min=3"`
	Body        string          `json:"body"`
	Blocks      json.RawMessage `json:"blocks" swaggertype:"object"`
	Type        string          `json:"type" validate:"required"`
	Attributes  string          `json:"attributes"`
	Status      string          `json:"status" validate:"required,oneof=DRAFT PUBLISHED SCHEDULED"`
	Language    string          `json:"language" validate:"required,len=2"`
	CategoryIDs []uint          `json:"category_ids"`
	Tags        []string        `json:"tags"` // Tag names
	PublishedAt *time.Time      `json:"published_at"`
}

type PaginatedContentResponse struct {
	Data []Content `json:"data"`
	Meta struct {
		Total int64 `json:"total"`
		Page  int   `json:"page"`
		Limit int   `json:"limit"`
	} `json:"meta"`
}

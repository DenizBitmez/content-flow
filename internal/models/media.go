package models

import "time"

type Media struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Filename  string    `json:"filename"`
	URL       string    `json:"url"`
	Size      int64     `json:"size"`
	ContentID uint      `json:"content_id"` // Optional link to content
	CreatedAt time.Time `json:"created_at"`
}

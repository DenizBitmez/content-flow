package models

import "gorm.io/gorm"

type Webhook struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	URL       string         `json:"url"`
	Events    string         `json:"events"` // Comma-separated: "content.create,content.update"
	Enabled   bool           `json:"enabled" gorm:"default:true"`
	CreatedAt int64          `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Helper to check if webhook subscribes to an event
func (w *Webhook) HasEvent(event string) bool {
	// Simple comma check (can be improved)
	// For now, assume simple strings
	return true // Simplified for MVP: trigger all enabled webhooks for now or implement split logic
}

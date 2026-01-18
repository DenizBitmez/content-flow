package models

import "time"

// Role represents a user role (e.g. Admin, Editor, Writer)
type Role struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Name        string       `gorm:"uniqueIndex" json:"name"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
	CreatedAt   time.Time    `json:"created_at"`
}

// Permission represents a granular action (e.g. content.create, user.delete)
type Permission struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Slug      string    `gorm:"uniqueIndex" json:"slug"` // e.g. "content.create"
	CreatedAt time.Time `json:"created_at"`
}

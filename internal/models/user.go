package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex" json:"username"`
	Email     string         `gorm:"uniqueIndex" json:"email"`
	FullName  string         `json:"full_name"`
	Bio       string         `json:"bio"`
	Avatar    string         `json:"avatar"`
	Password  string         `json:"-"`                        // Never return password in JSON
	RoleID    uint           `json:"role_id" gorm:"default:2"` // Default to Editor (2)
	Role      Role           `json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

package models

import "gorm.io/gorm"

type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex" json:"name"`
	Slug        string         `gorm:"uniqueIndex" json:"slug"`
	Description string         `json:"description"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Tag struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"uniqueIndex" json:"name"`
	Slug      string         `gorm:"uniqueIndex" json:"slug"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

package services

import (
	"content-flow/internal/database"
	"content-flow/internal/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

func CreateContent(content *models.Content) error {
	content.Version = 1
	return database.DB.Create(content).Error
}

func GetAllContent() ([]models.Content, error) {
	var contents []models.Content
	err := database.DB.Find(&contents).Error
	return contents, err
}

func GetContentByID(id uint) (*models.Content, error) {
	var content models.Content
	err := database.DB.First(&content, id).Error
	if err != nil {
		return nil, err
	}
	return &content, nil
}

// UpdateContent handles versioning: saves old state to ContentVersion, then updates Content
func UpdateContent(id uint, newTitle, newBody, newType, newAttributes, newStatus string) (*models.Content, error) {
	var content models.Content

	// Transaction guarantees atomicity
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Find existing content
		if err := tx.First(&content, id).Error; err != nil {
			return err
		}

		// 2. Create a snapshot (Version History)
		versionSnapshot := models.ContentVersion{
			ContentID:  content.ID,
			Title:      content.Title,
			Body:       content.Body,
			Type:       content.Type,
			Attributes: content.Attributes,
			Status:     content.Status,
			Version:    content.Version,
			ChangedAt:  time.Now(),
		}

		if err := tx.Create(&versionSnapshot).Error; err != nil {
			return err
		}

		// 3. Update the content and increment version
		content.Title = newTitle
		content.Body = newBody
		content.Type = newType
		content.Attributes = newAttributes
		content.Status = newStatus
		content.Version = content.Version + 1

		if err := tx.Save(&content).Error; err != nil {
			return err
		}

		return nil
	})

	return &content, err
}

func GetContentHistory(contentID uint) ([]models.ContentVersion, error) {
	var history []models.ContentVersion
	err := database.DB.Where("content_id = ?", contentID).Order("version desc").Find(&history).Error
	return history, err
}

func RevertContent(contentID uint, targetVersion int) (*models.Content, error) {
	var content models.Content
	var versionSnapshot models.ContentVersion

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Find current content
		if err := tx.First(&content, contentID).Error; err != nil {
			return err
		}

		// Find the target version
		if err := tx.Where("content_id = ? AND version = ?", contentID, targetVersion).First(&versionSnapshot).Error; err != nil {
			return errors.New("version not found")
		}

		// Save CURRENT state as history before reverting (so we don't lose the "bad" state)
		currentSnapshot := models.ContentVersion{
			ContentID:  content.ID,
			Title:      content.Title,
			Body:       content.Body,
			Type:       content.Type,
			Attributes: content.Attributes,
			Status:     content.Status,
			Version:    content.Version,
			ChangedAt:  time.Now(),
		}
		if err := tx.Create(&currentSnapshot).Error; err != nil {
			return err
		}

		// Revert Content to Snapshot Data
		// Note: We increment the version number to indicate a new change (the revert itself is a change)
		content.Title = versionSnapshot.Title
		content.Body = versionSnapshot.Body
		content.Type = versionSnapshot.Type
		content.Attributes = versionSnapshot.Attributes
		content.Status = versionSnapshot.Status
		content.Version = content.Version + 1

		if err := tx.Save(&content).Error; err != nil {
			return err
		}
		return nil
	})

	return &content, err
}

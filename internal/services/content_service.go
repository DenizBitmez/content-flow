package services

import (
	"content-flow/internal/database"
	"content-flow/internal/models"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func CreateContent(content *models.Content, categoryIDs []uint, tagNames []string, publishedAt *time.Time, blocks json.RawMessage, authorID uint) error {
	if content.Language == "" {
		content.Language = "en"
	}
	if content.GroupID == "" {
		content.GroupID = uuid.New().String()
	}
	content.AuthorID = authorID
	content.Version = 1

	if publishedAt != nil {
		content.PublishedAt = publishedAt
	}

	if len(blocks) > 0 {
		content.Blocks = datatypes.JSON(blocks)
	}

	// Handle Taxonomies
	if len(categoryIDs) > 0 {
		var categories []models.Category
		if err := database.DB.Where("id IN ?", categoryIDs).Find(&categories).Error; err != nil {
			return err
		}
		content.Categories = categories
	}

	if len(tagNames) > 0 {
		var tags []models.Tag
		for _, name := range tagNames {
			var tag models.Tag
			if err := database.DB.FirstOrCreate(&tag, models.Tag{Name: name, Slug: name}).Error; err != nil {
				return err
			}
			tags = append(tags, tag)
		}
		content.Tags = tags
	}

	if err := database.DB.Create(content).Error; err != nil {
		return err
	}

	// Trigger Webhook
	TriggerWebhooks("content.create", content)

	return nil
}

func AddTranslation(originalContentID uint, translation *models.Content) error {
	var original models.Content
	if err := database.DB.First(&original, originalContentID).Error; err != nil {
		return errors.New("original content not found")
	}

	// Verify if language already exists for this group
	var count int64
	database.DB.Model(&models.Content{}).Where("group_id = ? AND language = ?", original.GroupID, translation.Language).Count(&count)
	if count > 0 {
		return errors.New("translation for this language already exists")
	}

	translation.GroupID = original.GroupID
	translation.Version = 1
	// ID will be auto-generated because it's a new row
	return database.DB.Create(translation).Error
}

type ContentFilter struct {
	Search   string
	Type     string
	Status   string
	Language string
	Tags     []string
	Page     int
	Limit    int
}

func GetAllContent(filter ContentFilter) ([]models.Content, int64, error) {
	var contents []models.Content
	var total int64

	query := database.DB.Model(&models.Content{}).Preload("Categories").Preload("Tags")

	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where("title LIKE ? OR body LIKE ?", searchTerm, searchTerm)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Language != "" {
		query = query.Where("language = ?", filter.Language)
	}

	// Filtering by Tags (Join)
	if len(filter.Tags) > 0 {
		query = query.Joins("JOIN content_tags ON content_tags.content_id = contents.id").
			Joins("JOIN tags ON tags.id = content_tags.tag_id").
			Where("tags.name IN ?", filter.Tags).
			Group("contents.id") // Remove duplicates if multiple tags match
	}

	// Count total before pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	offset := (filter.Page - 1) * filter.Limit

	err := query.Limit(filter.Limit).Offset(offset).Order("created_at desc").Find(&contents).Error
	return contents, total, err
}

func GetContentByID(id uint) (*models.Content, error) {
	var content models.Content
	err := database.DB.Preload("Categories").Preload("Tags").First(&content, id).Error
	if err != nil {
		return nil, err
	}
	return &content, nil
}

// UpdateContent handles versioning: saves old state to ContentVersion, then updates Content
func UpdateContent(id uint, newTitle, newBody, newType, newAttributes, newStatus, newLang string, categoryIDs []uint, tagNames []string, publishedAt *time.Time, newBlocks json.RawMessage) (*models.Content, error) {
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
			Language:   content.Language,
			Blocks:     content.Blocks,
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
		if len(newBlocks) > 0 {
			content.Blocks = datatypes.JSON(newBlocks)
		}
		if newLang != "" {
			content.Language = newLang
		}
		if publishedAt != nil {
			content.PublishedAt = publishedAt
		}
		content.Version = content.Version + 1

		if err := tx.Save(&content).Error; err != nil {
			return err
		}

		// 4. Update Taxonomies
		// Categories
		if len(categoryIDs) > 0 {
			var categories []models.Category
			if err := tx.Where("id IN ?", categoryIDs).Find(&categories).Error; err != nil {
				return err
			}
			if err := tx.Model(&content).Association("Categories").Replace(categories); err != nil {
				return err
			}
		}

		// Tags
		if len(tagNames) > 0 {
			// Logic to sync tags inside transaction
			var tags []models.Tag
			for _, name := range tagNames {
				var tag models.Tag
				if err := tx.FirstOrCreate(&tag, models.Tag{Name: name, Slug: name}).Error; err != nil {
					return err
				}
				tags = append(tags, tag)
			}
			if err := tx.Model(&content).Association("Tags").Replace(tags); err != nil {
				return err
			}
		}

		return nil
	})

	// Trigger Webhook
	if err == nil {
		TriggerWebhooks("content.update", content)
	}

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
			Blocks:     content.Blocks,
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
		content.Blocks = versionSnapshot.Blocks
		content.Version = content.Version + 1

		if err := tx.Save(&content).Error; err != nil {
			return err
		}
		return nil
	})

	return &content, err
}

func PublishScheduledContent() {
	var scheduledContents []models.Content
	now := time.Now()

	// Find contents that are SCHEDULED and PublishedAt <= Now
	if err := database.DB.Where("status = ? AND published_at <= ?", "SCHEDULED", now).Find(&scheduledContents).Error; err != nil {
		return
	}

	for _, content := range scheduledContents {
		log.Printf("Publishing scheduled content ID: %d", content.ID)
		content.Status = "PUBLISHED"
		if err := database.DB.Save(&content).Error; err == nil {
			TriggerWebhooks("content.published", content)
		}
	}
}

func DeleteContent(id uint) error {
	// GORM soft delete
	if err := database.DB.Delete(&models.Content{}, id).Error; err != nil {
		return err
	}
	return nil
}

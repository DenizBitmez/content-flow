package services

import (
	"content-flow/internal/database"
	"content-flow/internal/models"
	"errors"

	"gorm.io/gorm"
)

func AddComment(userID, contentID uint, body string) (*models.Comment, error) {
	comment := models.Comment{
		UserID:    userID,
		ContentID: contentID,
		Body:      body,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		return nil, err
	}

	// Fetch with User info to return complete object
	if err := database.DB.Preload("User").First(&comment, comment.ID).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}

func GetComments(contentID uint) ([]models.Comment, error) {
	var comments []models.Comment
	if err := database.DB.Where("content_id = ?", contentID).Preload("User").Order("created_at desc").Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// ToggleLike adds a like if not exists, removes it if it does.
// Returns true if liked (added), false if unliked (removed).
func ToggleLike(userID, contentID uint) (bool, error) {
	var like models.Like
	result := database.DB.Where("user_id = ? AND content_id = ?", userID, contentID).First(&like)

	if result.Error == nil {
		// Like exists, remove it (Unlike)
		if err := database.DB.Delete(&like).Error; err != nil {
			return false, err
		}
		return false, nil
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Like does not exist, create it (Like)
		like = models.Like{
			UserID:    userID,
			ContentID: contentID,
		}
		if err := database.DB.Create(&like).Error; err != nil {
			return false, err
		}
		return true, nil
	} else {
		return false, result.Error
	}
}

func GetLikeCount(contentID uint) (int64, error) {
	var count int64
	if err := database.DB.Model(&models.Like{}).Where("content_id = ?", contentID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

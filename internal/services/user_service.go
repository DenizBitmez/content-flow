package services

import (
	"content-flow/internal/database"
	"content-flow/internal/models"
	"errors"
)

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func UpdateUserProfile(userID uint, bio, avatar, fullName string) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return err
	}

	user.Bio = bio
	user.Avatar = avatar
	user.FullName = fullName

	return database.DB.Save(&user).Error
}

func GetUserStories(authorID uint) ([]models.Content, error) {
	var stories []models.Content
	// Only fetch Published stories for public view? Or all for author?
	// For now let's just fetch all non-deleted.
	// Typically public profiles show only published.
	if err := database.DB.Where("author_id = ? AND status = 'PUBLISHED'", authorID).Order("created_at desc").Find(&stories).Error; err != nil {
		return nil, err
	}
	return stories, nil
}

package services

import (
	"content-flow/internal/database"
	"content-flow/internal/models"
)

// --- Categories ---

func CreateCategory(name, slug, description string) (*models.Category, error) {
	category := &models.Category{
		Name:        name,
		Slug:        slug,
		Description: description,
	}
	err := database.DB.Create(category).Error
	return category, err
}

func GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	err := database.DB.Find(&categories).Error
	return categories, err
}

func GetCategoriesByIDs(ids []uint) ([]models.Category, error) {
	var categories []models.Category
	if len(ids) == 0 {
		return categories, nil
	}
	err := database.DB.Where("id IN ?", ids).Find(&categories).Error
	return categories, err
}

// --- Tags ---

// SyncTags takes a list of tag names, finds existing ones, creates new ones, and returns the full list of Tag models.
func SyncTags(tagNames []string) ([]models.Tag, error) {
	var tags []models.Tag
	if len(tagNames) == 0 {
		return tags, nil
	}

	for _, name := range tagNames {
		var tag models.Tag
		// Find or Create
		// Note: Ideally use FirstOrCreate, but slug generation might be needed if complicated
		// For simplicity, assuming Slug = Name here or passed.
		// Let's assume Slug = Name for simple tags
		if err := database.DB.FirstOrCreate(&tag, models.Tag{Name: name, Slug: name}).Error; err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func GetAllTags() ([]models.Tag, error) {
	var tags []models.Tag
	err := database.DB.Find(&tags).Error
	return tags, err
}

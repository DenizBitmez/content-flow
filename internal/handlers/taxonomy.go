package handlers

import (
	"content-flow/internal/models"
	"content-flow/internal/pkgs/apierrors"
	"content-flow/internal/services"

	"github.com/gofiber/fiber/v2"
)

// CreateCategory godoc
// @Summary Create a new category
// @Tags Taxonomies
// @Accept json
// @Produce json
// @Param category body models.Category true "Category"
// @Success 200 {object} models.Category
// @Security Bearer
// @Router /api/categories [post]
func CreateCategory(c *fiber.Ctx) error {
	cat := new(models.Category)
	if err := c.BodyParser(cat); err != nil {
		return apierrors.BadRequest("Cannot parse JSON")
	}

	category, err := services.CreateCategory(cat.Name, cat.Slug, cat.Description)
	if err != nil {
		return apierrors.Internal(err.Error())
	}

	return c.JSON(category)
}

// GetAllCategories godoc
// @Summary Get all categories
// @Tags Taxonomies
// @Produce json
// @Success 200 {array} models.Category
// @Router /api/categories [get]
func GetAllCategories(c *fiber.Ctx) error {
	categories, err := services.GetAllCategories()
	if err != nil {
		return apierrors.Internal(err.Error())
	}
	return c.JSON(categories)
}

// GetAllTags godoc
// @Summary Get all tags
// @Tags Taxonomies
// @Produce json
// @Success 200 {array} models.Tag
// @Router /api/tags [get]
func GetAllTags(c *fiber.Ctx) error {
	tags, err := services.GetAllTags()
	if err != nil {
		return apierrors.Internal(err.Error())
	}
	return c.JSON(tags)
}

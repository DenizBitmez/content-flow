package handlers

import (
	"content-flow/internal/models"
	"content-flow/internal/pkgs/apierrors"
	"content-flow/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// CreateContent godoc
// @Summary Create new content
// @Description Creates a new content item
// @Tags Content
// @Accept json
// @Produce json
// @Param content body models.ContentCreateRequest true "Content object"
// @Success 200 {object} models.Content
// @Failure 400 {object} apierrors.AppError
// @Failure 500 {object} apierrors.AppError
// @Router /api/content [post]
func CreateContent(c *fiber.Ctx) error {
	req := new(models.ContentCreateRequest)
	if err := c.BodyParser(req); err != nil {
		return apierrors.BadRequest("Cannot parse JSON: " + err.Error())
	}

	content := &models.Content{
		Title:      req.Title,
		Slug:       req.Slug,
		Body:       req.Body,
		Type:       req.Type,
		Attributes: req.Attributes,
		Status:     req.Status,
	}

	if err := services.CreateContent(content); err != nil {
		return apierrors.Internal("Failed to create content: " + err.Error())
	}

	return c.JSON(content)
}

// GetAllContent godoc
// @Summary Get all content
// @Description Retrieves all content items
// @Tags Content
// @Produce json
// @Success 200 {array} models.Content
// @Failure 500 {object} apierrors.AppError
// @Router /api/content [get]
func GetAllContent(c *fiber.Ctx) error {
	contents, err := services.GetAllContent()
	if err != nil {
		return apierrors.Internal("Failed to retrieve contents: " + err.Error())
	}
	return c.JSON(contents)
}

// GetContent godoc
// @Summary Get content by ID
// @Description Retrieves a specific content item by ID
// @Tags Content
// @Produce json
// @Param id path int true "Content ID"
// @Success 200 {object} models.Content
// @Failure 404 {object} apierrors.AppError
// @Router /api/content/{id} [get]
func GetContent(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	content, err := services.GetContentByID(uint(id))
	if err != nil {
		return apierrors.NotFound("Content not found")
	}
	return c.JSON(content)
}

// UpdateContent godoc
// @Summary Update content
// @Description Updates an existing content item
// @Tags Content
// @Accept json
// @Produce json
// @Param id path int true "Content ID"
// @Param content body models.ContentUpdateRequest true "Update Request"
// @Success 200 {object} models.Content
// @Failure 400 {object} apierrors.AppError
// @Failure 500 {object} apierrors.AppError
// @Router /api/content/{id} [put]
func UpdateContent(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	req := new(models.ContentUpdateRequest)
	if err := c.BodyParser(req); err != nil {
		return apierrors.BadRequest("Cannot parse JSON: " + err.Error())
	}

	updatedContent, err := services.UpdateContent(uint(id), req.Title, req.Body, req.Type, req.Attributes, req.Status)
	if err != nil {
		return apierrors.Internal("Failed to update content: " + err.Error())
	}

	return c.JSON(updatedContent)
}

// GetHistory godoc
// @Summary Get content history
// @Description Retrieves version history for a content item
// @Tags Content
// @Produce json
// @Param id path int true "Content ID"
// @Success 200 {array} models.ContentVersion
// @Failure 500 {object} apierrors.AppError
// @Router /api/content/{id}/history [get]
func GetHistory(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	history, err := services.GetContentHistory(uint(id))
	if err != nil {
		return apierrors.Internal("Failed to retrieve history: " + err.Error())
	}
	return c.JSON(history)
}

// RevertContent godoc
// @Summary Revert content version
// @Description Reverts content to a specific version
// @Tags Content
// @Produce json
// @Param id path int true "Content ID"
// @Param version path int true "Version number"
// @Success 200 {object} models.Content
// @Failure 500 {object} apierrors.AppError
// @Router /api/content/{id}/revert/{version} [post]
func RevertContent(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	version, _ := strconv.Atoi(c.Params("version"))

	revertedContent, err := services.RevertContent(uint(id), version)
	if err != nil {
		return apierrors.Internal("Failed to revert content: " + err.Error())
	}
	return c.JSON(revertedContent)
}

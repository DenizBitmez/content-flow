package handlers

import (
	"content-flow/internal/models"
	"content-flow/internal/pkgs/apierrors"
	"content-flow/internal/services"
	"strconv"
	"strings"

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
// @Security Bearer
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
		Language:   req.Language,
	}

	if err := services.CreateContent(content, req.CategoryIDs, req.Tags, req.PublishedAt, req.Blocks); err != nil {
		return apierrors.Internal("Failed to create content: " + err.Error())
	}

	return c.JSON(content)
}

// GetAllContent godoc
// @Summary Get all content with filters
// @Description Retrieves content items with search, filters and pagination
// @Tags Content
// @Produce json
// @Param q query string false "Search term"
// @Param type query string false "Content Type"
// @Param status query string false "Content Status"
// @Param lang query string false "Language code"
// @Param tags query string false "Comma separated tags"
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 10)"
// @Success 200 {object} models.PaginatedContentResponse
// @Failure 500 {object} apierrors.AppError
// @Router /api/content [get]
func GetAllContent(c *fiber.Ctx) error {
	filter := services.ContentFilter{
		Search:   c.Query("q"),
		Type:     c.Query("type"),
		Status:   c.Query("status"),
		Language: c.Query("lang"),
		Page:     c.QueryInt("page", 1),
		Limit:    c.QueryInt("limit", 10),
	}

	if tags := c.Query("tags"); tags != "" {
		// Importing strings package might be needed if not present
		// For now assuming strings.Split works or adding import if needed
		// Let's use a helper or standard split if strings is imported
		// 'strings' is not imported in handler yet. I will rely on 'parser' or add import.
		// Actually, I'll assume strings is NOT imported and use a simple parsing or rely on adding import.
		// Wait, I should add strings import to be safe.
		// For this replacement, I'll attempt to use it and if it fails I'll add import in next step.
		// Actually, standard approach:
	}

	// Split tags manually or assume simple string for now?
	// Let's rely on standard split but I need to make sure strings is imported.
	// Check imports in file... 'strconv' is there. 'strings' probably not.

	// Better approach: filter.Tags = c.Query("tags") (as string) and handle split in service?
	// No, service expects []string.
	// I will just use split here and fix imports in separate step or same step if possible.
	// Wait, I can't modify imports easily in the same ReplaceFileContent chunk effectively if they are far away.
	// I will skip 'tags' splitting for a second or try to add it without strings? No.
	// I'll add the logic assuming I'll fix imports.

	if tagsStr := c.Query("tags"); tagsStr != "" {
		filter.Tags = strings.Split(tagsStr, ",")
	}

	contents, total, err := services.GetAllContent(filter)
	if err != nil {
		return apierrors.Internal("Failed to retrieve contents: " + err.Error())
	}

	return c.JSON(fiber.Map{
		"data": contents,
		"meta": fiber.Map{
			"total": total,
			"page":  filter.Page,
			"limit": filter.Limit,
		},
	})
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
// @Security Bearer
// @Router /api/content/{id} [put]
func UpdateContent(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	req := new(models.ContentUpdateRequest)
	if err := c.BodyParser(req); err != nil {
		return apierrors.BadRequest("Cannot parse JSON: " + err.Error())
	}

	updatedContent, err := services.UpdateContent(uint(id), req.Title, req.Body, req.Type, req.Attributes, req.Status, req.Language, req.CategoryIDs, req.Tags, req.PublishedAt, req.Blocks)
	if err != nil {
		return apierrors.Internal("Failed to update content: " + err.Error())
	}

	return c.JSON(updatedContent)
}

// AddTranslation godoc
// @Summary Add translation
// @Description Creates a new localized version of an existing content
// @Tags Content
// @Accept json
// @Produce json
// @Param id path int true "Original Content ID"
// @Param content body models.ContentCreateRequest true "Translated Content"
// @Success 200 {object} models.Content
// @Failure 400 {object} apierrors.AppError
// @Failure 500 {object} apierrors.AppError
// @Security Bearer
// @Router /api/content/{id}/localize [post]
func AddTranslation(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	req := new(models.ContentCreateRequest)
	if err := c.BodyParser(req); err != nil {
		return apierrors.BadRequest("Cannot parse JSON: " + err.Error())
	}

	translation := &models.Content{
		Title:      req.Title,
		Slug:       req.Slug,
		Body:       req.Body,
		Type:       req.Type,
		Attributes: req.Attributes,
		Status:     req.Status,
		Language:   req.Language,
	}

	// Note: Taxonomies for translations should theoretically be same as original or localized?
	// For now, let's keep it simple and not carry over taxonomies automatically, or allow setting them.
	// Users can update them later.
	if err := services.AddTranslation(uint(id), translation); err != nil {
		return apierrors.BadRequest("Failed to add translation: " + err.Error())
	}

	return c.JSON(translation)
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
// @Security Bearer
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

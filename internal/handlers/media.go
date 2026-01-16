package handlers

import (
	"content-flow/internal/database"
	"content-flow/internal/models"
	"content-flow/internal/pkgs/apierrors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UploadMedia godoc
// @Summary Upload media file
// @Description Uploads a media file (image) and associates it with optional content
// @Tags Media
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image file"
// @Param content_id formData int false "Content ID to associate"
// @Success 200 {object} models.Media
// @Failure 400 {object} apierrors.AppError
// @Failure 500 {object} apierrors.AppError
// @Security Bearer
// @Router /api/media [post]
func UploadMedia(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return apierrors.BadRequest("Image upload failed: " + err.Error())
	}

	// Generate a unique filename
	uniqueId := uuid.New()
	filename := fmt.Sprintf("%s-%s", uniqueId.String(), file.Filename)
	path := fmt.Sprintf("./uploads/%s", filename)

	if err := c.SaveFile(file, path); err != nil {
		return apierrors.Internal("Failed to save file: " + err.Error())
	}

	// Save to DB
	media := models.Media{
		Filename:  filename,
		URL:       "/uploads/" + filename,
		Size:      file.Size,
		CreatedAt: time.Now(),
	}

	if contentID := c.FormValue("content_id"); contentID != "" {
		var cid uint
		if _, err := fmt.Sscanf(contentID, "%d", &cid); err == nil {
			media.ContentID = cid
		}
	}

	if result := database.DB.Create(&media); result.Error != nil {
		return apierrors.Internal("Database error: " + result.Error.Error())
	}

	return c.JSON(media)
}

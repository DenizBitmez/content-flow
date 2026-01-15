package handlers

import (
	"content-flow/internal/database"
	"content-flow/internal/models"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UploadMedia(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Image upload failed"})
	}

	// Generate a unique filename
	uniqueId := uuid.New()
	filename := fmt.Sprintf("%s-%s", uniqueId.String(), file.Filename)
	path := fmt.Sprintf("./uploads/%s", filename)

	if err := c.SaveFile(file, path); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save file"})
	}

	// Save to DB
	media := models.Media{
		Filename:  filename,
		URL:       "/uploads/" + filename,
		Size:      file.Size,
		CreatedAt: time.Now(),
	}

	if contentID := c.FormValue("content_id"); contentID != "" {
		// Simple conversion, ignoring error for brevity in this step, ideally should handle
		var cid uint
		fmt.Sscanf(contentID, "%d", &cid)
		media.ContentID = cid
	}

	database.DB.Create(&media)

	return c.JSON(media)
}

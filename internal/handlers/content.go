package handlers

import (
	"content-flow/internal/models"
	"content-flow/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreateContent(c *fiber.Ctx) error {
	content := new(models.Content)
	if err := c.BodyParser(content); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := services.CreateContent(content); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(content)
}

func GetAllContent(c *fiber.Ctx) error {
	contents, err := services.GetAllContent()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(contents)
}

func GetContent(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	content, err := services.GetContentByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Content not found"})
	}
	return c.JSON(content)
}

func UpdateContent(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	type UpdateRequest struct {
		Title      string `json:"title"`
		Body       string `json:"body"`
		Type       string `json:"type"`
		Attributes string `json:"attributes"`
		Status     string `json:"status"`
	}
	req := new(UpdateRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	updatedContent, err := services.UpdateContent(uint(id), req.Title, req.Body, req.Type, req.Attributes, req.Status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(updatedContent)
}

func GetHistory(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	history, err := services.GetContentHistory(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(history)
}

func RevertContent(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	version, _ := strconv.Atoi(c.Params("version"))

	revertedContent, err := services.RevertContent(uint(id), version)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(revertedContent)
}

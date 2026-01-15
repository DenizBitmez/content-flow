package main

import (
	"content-flow/internal/database"
	"content-flow/internal/handlers"
	"content-flow/internal/models"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// 1. Connect to Database
	database.Connect()

	// 2. Run Auto-Migrations
	log.Println("Running Auto-migrations...")
	database.DB.AutoMigrate(&models.Content{}, &models.ContentVersion{}, &models.Media{})

	// 3. Setup Fiber App
	app := fiber.New()
	app.Use(logger.New())

	// Static route for uploads
	app.Static("/uploads", "./uploads")

	// 4. Setup Routes
	api := app.Group("/api")

	// Media
	api.Post("/media", handlers.UploadMedia)

	// Content
	api.Post("/content", handlers.CreateContent)
	api.Get("/content", handlers.GetAllContent)
	api.Get("/content/:id", handlers.GetContent)
	api.Put("/content/:id", handlers.UpdateContent)

	// History / Versioning
	api.Get("/content/:id/history", handlers.GetHistory)
	api.Post("/content/:id/revert/:version", handlers.RevertContent)

	// 5. Start Server
	log.Fatal(app.Listen(":3000"))
}

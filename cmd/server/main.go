package main

import (
	"content-flow/internal/database"
	"content-flow/internal/handlers"
	"content-flow/internal/models"
	"content-flow/internal/pkgs/apierrors"
	"log"

	_ "content-flow/docs" // Import generated swagger docs

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

// @title ContentFlow CMS API
// @version 1.0
// @description Headless CMS API with dynamic content and media management
// @contact.name API Support
// @contact.email support@contentflow.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:3000
// @BasePath /
func main() {
	// 1. Connect to Database
	database.Connect()

	// 2. Run Auto-Migrations
	log.Println("Running Auto-migrations...")
	database.DB.AutoMigrate(&models.Content{}, &models.ContentVersion{}, &models.Media{})

	// 3. Setup Fiber App with Global Error Handler
	app := fiber.New(fiber.Config{
		ErrorHandler: apierrors.ErrorHandler,
	})
	app.Use(logger.New())

	// Static route for uploads
	app.Static("/uploads", "./uploads")

	// Swagger Middleware
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/swagger/doc.json",
		DeepLinking: false,
	}))

	// Redirect /swagger to /swagger/index.html
	app.Get("/swagger", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html", fiber.StatusMovedPermanently)
	})

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

package main

import (
	"content-flow/internal/database"
	"content-flow/internal/handlers"
	"content-flow/internal/models"
	"content-flow/internal/pkgs/apierrors"
	"content-flow/internal/pkgs/auth"
	"content-flow/internal/services"
	"log"
	"time"

	_ "content-flow/docs" // Import generated swagger docs

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
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
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// 1. Connect to Database
	database.Connect()

	// 2. Run Auto-Migrations
	log.Println("Running Auto-migrations...")
	database.DB.AutoMigrate(&models.Content{}, &models.ContentVersion{}, &models.Media{}, &models.User{}, &models.Category{}, &models.Tag{}, &models.Webhook{})

	// 3. Setup Fiber App with Global Error Handler
	app := fiber.New(fiber.Config{
		ErrorHandler: apierrors.ErrorHandler,
	})
	app.Use(logger.New())

	// General Rate Limiter (60 req/min)
	app.Use(limiter.New(limiter.Config{
		Max:        60,
		Expiration: 1 * time.Minute,
	}))

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

	// Public Taxonomy Routes
	api.Get("/categories", handlers.GetAllCategories)
	api.Get("/tags", handlers.GetAllTags)

	// Auth Rate Limiter (5 req/min) - Brute Force Protection
	authLimiter := limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return apierrors.New(429, "Too many login attempts. Please wait a minute.")
		},
	})

	// Auth Routes
	authGroup := api.Group("/auth", authLimiter)
	authGroup.Post("/register", handlers.Register)
	authGroup.Post("/login", handlers.Login)

	// Media Rate Limiter (5 req/min)
	uploadLimiter := limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return apierrors.New(429, "Too many upload requests. Please wait a minute.")
		},
	})

	// Private Routes (Require Token)
	private := api.Group("/", auth.Protected())

	// Taxonomies
	private.Post("/categories", handlers.CreateCategory)

	// Webhooks
	private.Post("/webhooks", handlers.CreateWebhook)
	private.Get("/webhooks", handlers.GetAllWebhooks)

	// Media
	private.Post("/media", uploadLimiter, handlers.UploadMedia)

	// Content
	private.Post("/content", handlers.CreateContent)
	private.Post("/content/:id/localize", handlers.AddTranslation)
	// Public Read Access for Content
	api.Get("/content", handlers.GetAllContent)
	api.Get("/content/:id", handlers.GetContent)
	// Private Update
	private.Put("/content/:id", handlers.UpdateContent)

	// History / Versioning
	private.Get("/content/:id/history", handlers.GetHistory)
	private.Post("/content/:id/revert/:version", handlers.RevertContent)

	// 5. Start Scheduler (Background Job)
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			services.PublishScheduledContent()
		}
	}()

	// 6. Start Server
	log.Fatal(app.Listen(":3000"))
}

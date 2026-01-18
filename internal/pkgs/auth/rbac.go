package auth

import (
	"content-flow/internal/database"
	"content-flow/internal/models"
	"content-flow/internal/pkgs/apierrors"

	"github.com/gofiber/fiber/v2"
)

// RequirePermission checks if the authenticated user has the specified permission
func RequirePermission(permSlug string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		if userID == nil {
			return apierrors.New(fiber.StatusUnauthorized, "Unauthorized")
		}

		// Fetch User with Role and Permissions
		var user models.User
		if err := database.DB.Preload("Role.Permissions").First(&user, userID).Error; err != nil {
			return apierrors.New(fiber.StatusUnauthorized, "User not found")
		}

		// Check if user has admin role (bypass check)
		if user.Role.Name == "Admin" {
			return c.Next()
		}

		// Check if role has the permission
		hasPerm := false
		for _, p := range user.Role.Permissions {
			if p.Slug == permSlug {
				hasPerm = true
				break
			}
		}

		if !hasPerm {
			return apierrors.New(fiber.StatusForbidden, "Forbidden: Missing permission "+permSlug)
		}

		return c.Next()
	}
}

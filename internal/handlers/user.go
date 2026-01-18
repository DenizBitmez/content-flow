package handlers

import (
	"content-flow/internal/pkgs/apierrors"
	"content-flow/internal/services"

	"github.com/gofiber/fiber/v2"
)

type UpdateProfileRequest struct {
	Bio      string `json:"bio"`
	Avatar   string `json:"avatar"`
	FullName string `json:"full_name"`
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get public profile by username
// @Tags Users
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} models.User
// @Failure 404 {object} apierrors.AppError
// @Router /api/users/{username} [get]
func GetProfile(c *fiber.Ctx) error {
	username := c.Params("username")
	user, err := services.GetUserByUsername(username)
	if err != nil {
		return apierrors.NotFound("User not found")
	}
	return c.JSON(user)
}

// GetUserStories godoc
// @Summary Get user stories
// @Description Get published stories by username
// @Tags Users
// @Produce json
// @Param username path string true "Username"
// @Success 200 {array} models.Content
// @Router /api/users/{username}/stories [get]
func GetUserStories(c *fiber.Ctx) error {
	username := c.Params("username")
	user, err := services.GetUserByUsername(username)
	if err != nil {
		return apierrors.NotFound("User not found")
	}

	stories, err := services.GetUserStories(user.ID)
	if err != nil {
		return apierrors.Internal(err.Error())
	}

	return c.JSON(stories)
}

// UpdateProfile godoc
// @Summary Update profile
// @Description Update bio, avatar, full name
// @Tags Users
// @Accept json
// @Produce json
// @Param request body UpdateProfileRequest true "Profile Update"
// @Success 200 {object} map[string]string
// @Security Bearer
// @Router /api/users/profile [put]
func UpdateProfile(c *fiber.Ctx) error {
	userID := uint(c.Locals("user_id").(float64))
	req := new(UpdateProfileRequest)
	if err := c.BodyParser(req); err != nil {
		return apierrors.BadRequest("Cannot parse JSON")
	}

	if err := services.UpdateUserProfile(userID, req.Bio, req.Avatar, req.FullName); err != nil {
		return apierrors.Internal(err.Error())
	}

	return c.JSON(fiber.Map{"message": "Profile updated successfully"})
}

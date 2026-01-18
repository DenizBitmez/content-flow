package handlers

import (
	"content-flow/internal/pkgs/apierrors"
	"content-flow/internal/pkgs/validator"
	"content-flow/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Token   string `json:"token,omitempty"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Username string `json:"username" validate:"required,min=3,alphanum"`
	FullName string `json:"full_name" validate:"required"`
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register Request"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} apierrors.AppError
// @Router /api/auth/register [post]
func Register(c *fiber.Ctx) error {
	req := new(RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return apierrors.BadRequest("Cannot parse JSON")
	}

	if errors := validator.ValidateStruct(req); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  errors,
			"message": "Validation failed",
		})
	}

	if err := services.RegisterUser(req.Email, req.Password, req.Username, req.FullName); err != nil {
		return apierrors.BadRequest(err.Error())
	}

	return c.JSON(AuthResponse{
		Success: true,
		Message: "User registered successfully",
	})
}

// Login godoc
// @Summary Login user
// @Description Logs in a user and returns a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body AuthRequest true "Login Request"
// @Success 200 {object} AuthResponse
// @Failure 401 {object} apierrors.AppError
// @Router /api/auth/login [post]
func Login(c *fiber.Ctx) error {
	req := new(AuthRequest)
	if err := c.BodyParser(req); err != nil {
		return apierrors.BadRequest("Cannot parse JSON")
	}

	token, err := services.LoginUser(req.Email, req.Password)
	if err != nil {
		return apierrors.New(fiber.StatusUnauthorized, "Invalid credentials")
	}

	return c.JSON(AuthResponse{
		Success: true,
		Token:   token,
	})
}

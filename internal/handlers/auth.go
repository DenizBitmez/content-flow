package handlers

import (
	"content-flow/internal/pkgs/apierrors"
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

// Register godoc
// @Summary Register a new user
// @Description Creates a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body AuthRequest true "Register Request"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} apierrors.AppError
// @Router /api/auth/register [post]
func Register(c *fiber.Ctx) error {
	req := new(AuthRequest)
	if err := c.BodyParser(req); err != nil {
		return apierrors.BadRequest("Cannot parse JSON")
	}

	if err := services.RegisterUser(req.Email, req.Password); err != nil {
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

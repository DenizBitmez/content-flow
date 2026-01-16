package apierrors

import "github.com/gofiber/fiber/v2"

// AppError represents a standard API error
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// New creates a new AppError
func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Success: false,
	}
}

// BadRequest creates a 400 Bad Request error
func BadRequest(message string) *AppError {
	return New(fiber.StatusBadRequest, message)
}

// NotFound creates a 404 Not Found error
func NotFound(message string) *AppError {
	return New(fiber.StatusNotFound, message)
}

// Internal creates a 500 Internal Server Error
func Internal(message string) *AppError {
	return New(fiber.StatusInternalServerError, message)
}

// ErrorHandler is the global error handler for Fiber
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default error code
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	// Check if it's our custom AppError
	if e, ok := err.(*AppError); ok {
		code = e.Code
		message = e.Message
	} else if e, ok := err.(*fiber.Error); ok {
		// Handle Fiber's built-in errors
		code = e.Code
		message = e.Message
	} else if err != nil {
		// Handle standard Go errors
		message = err.Error()
	}

	// Make sure we don't leak specialized internal errors if we don't want to
	// (For now, we just pass the message, but in production we might want to sanitize 500s)

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": message,
		"code":    code,
	})
}

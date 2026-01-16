package auth

import (
	"content-flow/internal/pkgs/apierrors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte("super-secret-key-change-this-in-env")

func GenerateToken(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return apierrors.New(fiber.StatusUnauthorized, "Missing Authorization Header")
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, apierrors.New(fiber.StatusUnauthorized, "Unexpected signing method")
			}
			return SecretKey, nil
		})

		if err != nil || !token.Valid {
			return apierrors.New(fiber.StatusUnauthorized, "Invalid or Expired Token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return apierrors.Internal("Invalid Token Claims")
		}

		// Store user info in locals
		c.Locals("user_id", claims["user_id"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}

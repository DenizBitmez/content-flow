package services

import (
	"content-flow/internal/database"
	"content-flow/internal/models"
	"content-flow/internal/pkgs/auth"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(email, password, username, fullName string) error {
	// Check if email exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return errors.New("email already exists")
	}

	// Check if username exists
	if err := database.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return errors.New("username already taken")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Fetch Default Role (Editor)
	var role models.Role
	if err := database.DB.Where("name = ?", "Editor").First(&role).Error; err != nil {
		// Fallback to ID 2 if not found or handled by seeder
		role.ID = 2
	}

	user := models.User{
		Email:    email,
		Password: string(hashedPassword),
		Username: username,
		FullName: fullName,
		RoleID:   role.ID,
	}

	return database.DB.Create(&user).Error
}

func LoginUser(email, password string) (string, error) {
	var user models.User
	if err := database.DB.Preload("Role").Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := auth.GenerateToken(user.ID, user.Role.Name)
	if err != nil {
		return "", err
	}

	return token, nil
}

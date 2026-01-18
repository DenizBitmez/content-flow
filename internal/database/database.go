package database

import (
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default/env values")
	}

	var err error
	driver := os.Getenv("DB_DRIVER")
	dsn := os.Getenv("DB_DSN")

	if driver == "postgres" {
		log.Println("Connecting to PostgreSQL...")
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	} else {
		log.Println("Connecting to SQLite...")
		DB, err = gorm.Open(sqlite.Open("content_flow.db"), &gorm.Config{})
	}

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Connected to Database")
}

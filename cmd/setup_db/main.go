package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Connect to default 'postgres' database
	dsn := "host=localhost user=postgres password=12345 dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Could not connect to 'postgres' database. Check credentials/running status:", err)
	}

	// 2. Check if 'contentflow' exists
	var count int64
	db.Raw("SELECT count(*) FROM pg_database WHERE datname = ?", "contentflow").Scan(&count)

	if count == 0 {
		// 3. Create Database
		// Note: CREATE DATABASE cannot run inside a transaction, so we use Exec directly on the underlying connection if possible,
		// but GORM Exec should work if not in transaction.
		if err := db.Exec("CREATE DATABASE contentflow").Error; err != nil {
			log.Fatal("Failed to create database:", err)
		}
		fmt.Println("✅ Database 'contentflow' created successfully!")
	} else {
		fmt.Println("ℹ️ Database 'contentflow' already exists.")
	}
}

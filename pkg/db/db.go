// /pkg/db/db.go
package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, *sql.DB) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		QueryFields: true,
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // log to stdout
			logger.Config{
				SlowThreshold:             time.Second, // log slow queries
				LogLevel:                  logger.Info, // Log level (Info will log all queries)
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound when logging
				Colorful:                  true,        // Enable colorful output
			},
		),
	})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Retrieve the underlying *sql.DB object from GORM
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get the underlying database connection: %v", err)
	}

	// Set the maximum number of idle connections in the pool
	sqlDB.SetMaxIdleConns(10)

	// Set the maximum number of open connections to the database
	sqlDB.SetMaxOpenConns(100)

	// Set the maximum amount of time a connection may be reused (connection lifetime)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection established with connection pooling settings")

	return DB, sqlDB
}

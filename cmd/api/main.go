package main

import (
	"context"
	"inventory_management/api/handler"
	"inventory_management/internal/repository"
	"inventory_management/internal/usecase"
	"inventory_management/pkg/db"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"net/http"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Read ReadHeaderTimeout from the .env file and convert to time.Duration
	readHeaderTimeoutStr := os.Getenv("READ_HEADER_TIMEOUT")
	readHeaderTimeout, err := strconv.Atoi(readHeaderTimeoutStr)
	if err != nil || readHeaderTimeout <= 0 {
		log.Warn("Invalid or missing READ_HEADER_TIMEOUT, defaulting to 10 seconds")
		readHeaderTimeout = 10 // default to 10 seconds if the env variable is invalid or missing
	}

	// Initialize DB connection
	db, sqlDB := db.InitDB(false)
	if db == nil || sqlDB == nil {
		log.Fatal("Failed to initialize the database.")
	}

	// Initialize repository, use case, and handler
	productRepo := repository.NewPostgresProductRepository(db)
	productUsecase := usecase.NewProductUsecase(productRepo)
	productHandler := handler.NewProductHandler(productUsecase)

	// Setup the router by calling the new SetupRouter function
	router := SetupRouter(productHandler)

	// Create the HTTP server with the Gin router as its handler
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: time.Duration(readHeaderTimeout) * time.Second, // Adding ReadHeaderTimeout to prevent Slowloris attack
	}

	// Start the server in a goroutine so that it doesn't block graceful shutdown handling
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("Server running on port 8080")

	// Create a channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // Listen for SIGINT and SIGTERM

	// Block until a signal is received
	sig := <-quit
	log.WithFields(log.Fields{
		"signal": sig,
	}).Println("Received shutdown signal, shutting down server...")

	// Create a context with a timeout to allow for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Close the database connection
	if err := sqlDB.Close(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to close database connection")
	}

	log.Println("Server and database connection closed gracefully")

	// Attempt graceful shutdown by stopping the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

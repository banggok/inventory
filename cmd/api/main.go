// /cmd/api/main.go
package main

import (
	"context"
	"inventory_management/api/handler"
	product_repository "inventory_management/internal/repository"
	"inventory_management/internal/usecase"
	"inventory_management/pkg/db"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize DB connection
	db, sqlDB := db.InitDB()

	// Initialize repository, usecase and handler
	productRepo := product_repository.NewPostgresProductRepository(db)
	productUsecase := usecase.NewProductUsecase(productRepo)
	productHandler := handler.NewProductHandler(productUsecase)

	// Setup Gin
	router := gin.Default()

	// Define Routes
	router.POST("/products", productHandler.CreateProduct)
	router.GET("/products/:id", productHandler.GetProduct)

	// Create the http.Server with the Gin router as its handler
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start the server in a goroutine so that it doesn't block the graceful shutdown handling
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
	<-quit
	log.Println("Shutting down server...")

	// Create a context with a timeout to allow for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Close the database connection
	if err := sqlDB.Close(); err != nil {
		log.Fatalf("Failed to close database connection: %v", err)
	}

	log.Println("Server and database connection closed gracefully")

	// Attempt graceful shutdown by stopping the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

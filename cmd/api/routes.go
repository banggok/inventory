// /cmd/api/routes.go
package main

import (
	"inventory_management/api/handler"

	"github.com/gin-gonic/gin"
)

// SetupRouter defines all the application routes and returns the Gin router
func SetupRouter(productHandler *handler.ProductHandler) *gin.Engine {
	router := gin.Default()

	// Define Routes with route grouping
	api := router.Group("/api/v1")
	{
		api.POST("/products", productHandler.CreateProduct)
		api.GET("/products/:id", productHandler.GetProduct)
		api.PUT("/products/:id", productHandler.UpdateProductName) // Add the route for updating the product name

	}

	return router
}

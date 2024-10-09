// /api/handler/product_handler.go
package handler

import (
	"inventory_management/api/handler/dto" // Import the DTO package
	"inventory_management/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productUsecase usecase.ProductUsecase
}

func NewProductHandler(u usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{productUsecase: u}
}

// CreateProduct handles the creation of a new product
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest // Use the DTO for the request

	// Bind the request body to the req struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the request using the Validate method in the DTO
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call usecase to create a new product
	product, err := h.productUsecase.CreateProduct(req) // No SKU passed
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the created product
	c.JSON(http.StatusCreated, product)
}

// GetProduct retrieves a product by its ID
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Call usecase to get product by ID
	product, err := h.productUsecase.GetProductByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	// Respond with the product
	c.JSON(http.StatusOK, product)
}

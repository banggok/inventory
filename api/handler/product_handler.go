// /api/handler/product_handler.go
package handler

import (
	"inventory_management/api/handler/dto"
	transformer "inventory_management/api/transform"
	"inventory_management/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type ProductHandler struct {
	productUsecase usecase.ProductUsecase
}

func NewProductHandler(u usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{productUsecase: u}
}

// CreateProduct handles the creation of a new product
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest

	// Bind the request body to the req struct and validate it
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Invalid product creation request")
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Call usecase to create a new product using the request data
	product, err := h.productUsecase.CreateProduct(req.Name)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"name":  req.Name,
		}).Error("Error creating product")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	// Transform the entity to a response DTO
	productResponse := transformer.TransformProductEntityToResponse(product)

	// Respond with the created product
	log.WithFields(log.Fields{
		"id":   product.ID(),
		"name": product.Name(),
	}).Info("Product created successfully")
	c.JSON(http.StatusCreated, productResponse)
}

// GetProduct retrieves a product by its ID
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || id == 0 {
		log.WithFields(log.Fields{
			"error": err,
			"id":    idParam,
		}).Warn("Invalid product ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Call usecase to get product by ID
	product, err := h.productUsecase.GetProductByID(uint(id))
	if err != nil {
		// If the error is "product not found", return 404
		if err.Error() == "product not found" {
			log.WithFields(log.Fields{
				"id": id,
			}).Warn("Product not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			// Log any other errors and return 500 Internal Server Error
			log.WithFields(log.Fields{
				"error": err,
				"id":    id,
			}).Error("Error retrieving product from the database")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Transform the entity to a response DTO
	productResponse := transformer.TransformProductEntityToResponse(product)

	// Respond with the product
	log.WithFields(log.Fields{
		"id":   product.ID(),
		"name": product.Name(),
	}).Info("Product retrieved successfully")
	c.JSON(http.StatusOK, productResponse)
}

package handler

import (
	"inventory_management/api/handler/dto"
	helper_handler "inventory_management/api/handler/helper"
	"inventory_management/api/handler/transformer"
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
	validationErrors, err := helper_handler.ReadAndValidateRequestBody(c, &req)
	if validationErrors != nil {
		helper_handler.SendErrorResponse(c, http.StatusUnprocessableEntity, validationErrors)
		return
	}

	if err != nil {
		helper_handler.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Call use case to create a new product
	product, err := h.productUsecase.CreateProduct(req.Name)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"name":  req.Name,
		}).Error("Failed to create product")
		helper_handler.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create product")
		return
	}

	// Transform the entity to response DTO and respond
	productResponse := transformer.TransformProductEntityToResponse(product)
	log.WithFields(log.Fields{"id": product.ID(), "name": product.Name()}).Info("Product created successfully")
	c.JSON(http.StatusCreated, productResponse)
}

// GetProduct retrieves a product by its ID
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || id == 0 {
		log.Warnf("Invalid product ID: %v", idParam)
		helper_handler.SendErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.productUsecase.GetProductByID(uint(id))
	if err != nil {
		if err.Error() == "product not found" {
			log.Warnf("Product not found: %d", id)
			helper_handler.SendErrorResponse(c, http.StatusNotFound, "Product not found")
		} else {
			log.Errorf("Failed to retrieve product: %d, error: %v", id, err)
			helper_handler.SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve product")
		}
		return
	}

	// Transform the entity to response DTO and respond
	productResponse := transformer.TransformProductEntityToResponse(product)
	log.Infof("Product retrieved successfully: ID: %d", product.ID())
	c.JSON(http.StatusOK, productResponse)
}

// UpdateProductName handles updating a product's name
func (h *ProductHandler) UpdateProductName(c *gin.Context) {
	var req dto.UpdateProductRequest
	validationErrors, err := helper_handler.ReadAndValidateRequestBody(c, &req)
	if validationErrors != nil {
		helper_handler.SendErrorResponse(c, http.StatusUnprocessableEntity, validationErrors)
		return
	}

	if err != nil {
		helper_handler.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Parse the product ID from the URL
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || id == 0 {
		helper_handler.SendErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	// Call use case to update the product name
	product, err := h.productUsecase.UpdateProductName(uint(id), req.Name)
	if err != nil {
		log.Errorf("Failed to update product name: %d, error: %v", id, err)
		if err.Error() == "product not found" {
			helper_handler.SendErrorResponse(c, http.StatusNotFound, "Product not found")
		} else {
			helper_handler.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update product")
		}
		return
	}

	// Transform the entity to response DTO and respond
	productResponse := transformer.TransformProductEntityToResponse(product)
	log.Infof("Product updated successfully: ID: %d", product.ID())
	c.JSON(http.StatusOK, productResponse)
}

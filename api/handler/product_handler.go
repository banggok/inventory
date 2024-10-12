package handler

import (
	consts "inventory_management/api/handler/const"
	"inventory_management/api/handler/dto"
	helper_handler "inventory_management/api/handler/helper"
	"inventory_management/api/handler/transformer"
	"inventory_management/internal/usecase"
	"inventory_management/pkg/utility"
	"net/http"

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
	var req dto.CreateProductRequest

	// Read and validate the request body
	validationErrors, err := helper_handler.ReadAndValidateRequestBody(c, &req)
	if validationErrors != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": validationErrors})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// Create product using the usecase
	product, err := h.productUsecase.CreateProduct(req.Name)
	if err != nil {
		// If there's an error creating the product, return the correct error message
		helper_handler.HandleErrorResponse(c, err, consts.ErrFailedCreate, http.StatusInternalServerError)
		return
	}

	// Transform and send a success response
	productResponse := transformer.TransformProductEntityToResponse(product)
	utility.LogSuccess("product created successfully", product.ID(), product.Name())
	c.JSON(http.StatusCreated, productResponse)
}

// GetProduct retrieves a product by its ID
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := helper_handler.ParseIDFromParam(c)
	if err != nil {
		helper_handler.HandleErrorResponse(c, err, consts.ErrInvalidProductID, http.StatusBadRequest)
		return
	}

	product, err := h.productUsecase.GetProductByID(id)
	if err != nil {
		if err == usecase.ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{"errors": consts.ErrProductNotFound})
		} else {
			helper_handler.HandleErrorResponse(c, err, consts.ErrFailedRetrieve, http.StatusInternalServerError)
		}
		return
	}

	productResponse := transformer.TransformProductEntityToResponse(product)
	utility.LogSuccess("product retrieved successfully", product.ID(), product.Name())
	c.JSON(http.StatusOK, productResponse)
}

// UpdateProductName handles updating a product's name
func (h *ProductHandler) UpdateProductName(c *gin.Context) {
	var req dto.UpdateProductRequest

	validationErrors, err := helper_handler.ReadAndValidateRequestBody(c, &req)
	if validationErrors != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": validationErrors})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	id, err := helper_handler.ParseIDFromParam(c)
	if err != nil {
		helper_handler.HandleErrorResponse(c, err, consts.ErrInvalidProductID, http.StatusBadRequest)
		return
	}

	product, err := h.productUsecase.UpdateProductName(id, req.Name)
	if err != nil {
		if err == usecase.ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{"errors": consts.ErrProductNotFound})
		} else {
			helper_handler.HandleErrorResponse(c, err, consts.ErrFailedUpdate, http.StatusInternalServerError)
		}
		return
	}

	productResponse := transformer.TransformProductEntityToResponse(product)
	utility.LogSuccess("product updated successfully", product.ID(), product.Name())
	c.JSON(http.StatusOK, productResponse)
}

// GetProductList handles listing products with filters, sorting, and pagination
func (h *ProductHandler) GetProductList(c *gin.Context) {
	queryParams := dto.ProductListQueryParams{}

	// Perform validation manually
	validationErrors := queryParams.Validate(c.Request.URL.Query())
	if validationErrors != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": validationErrors})
		return
	}

	// Fetch the products based on filters, sorting, and pagination
	products, err := h.productUsecase.ListProducts(
		queryParams.SearchTerm,
		queryParams.SortBy,
		queryParams.SortDirection,
		queryParams.Limit,
		queryParams.Offset,
	)
	if err != nil {
		helper_handler.HandleErrorResponse(c, err, consts.ErrFailedRetrieve, http.StatusInternalServerError)
		return
	}

	// Transform the products to response DTOs
	productResponses := make([]*dto.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = transformer.TransformProductEntityToResponse(product)
	}

	utility.LogSuccess("product list retrieved successfully", len(products), "products")
	c.JSON(http.StatusOK, gin.H{
		"products": productResponses,
		"total":    len(products),
	})
}

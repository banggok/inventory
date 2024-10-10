package handler

import (
	consts "inventory_management/api/handler/const"
	"inventory_management/api/handler/dto"
	helper_handler "inventory_management/api/handler/helper"
	"inventory_management/api/handler/transformer"
	helper "inventory_management/helper" // Import global helper for logging
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

	product, err := h.productUsecase.CreateProduct(req.Name)
	if err != nil {
		helper.LogError(consts.ErrFailedCreate, req.Name, err) // Use global helper
		helper_handler.SendErrorResponse(c, http.StatusInternalServerError, consts.ErrFailedCreate)
		return
	}

	productResponse := transformer.TransformProductEntityToResponse(product)
	helper.LogSuccess("Product created", product.ID(), product.Name()) // Use global helper
	c.JSON(http.StatusCreated, productResponse)
}

// GetProduct retrieves a product by its ID
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || id == 0 {
		helper.LogError(consts.ErrInvalidProductID, idParam, err) // Use global helper
		helper_handler.SendErrorResponse(c, http.StatusBadRequest, consts.ErrInvalidProductID)
		return
	}

	product, err := h.productUsecase.GetProductByID(uint(id))
	if err != nil {
		helper.LogError(consts.ErrFailedRetrieve, strconv.Itoa(int(id)), err) // Use global helper
		if err.Error() == "product not found" {
			helper_handler.SendErrorResponse(c, http.StatusNotFound, consts.ErrProductNotFound)
		} else {
			helper_handler.SendErrorResponse(c, http.StatusInternalServerError, consts.ErrFailedRetrieve)
		}
		return
	}

	productResponse := transformer.TransformProductEntityToResponse(product)
	helper.LogSuccess("Product retrieved", product.ID(), product.Name()) // Use global helper
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

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || id == 0 {
		helper_handler.SendErrorResponse(c, http.StatusBadRequest, consts.ErrInvalidProductID)
		return
	}

	product, err := h.productUsecase.UpdateProductName(uint(id), req.Name)
	if err != nil {
		helper.LogError(consts.ErrFailedUpdate, strconv.Itoa(int(id)), err) // Use global helper
		if err.Error() == "product not found" {
			helper_handler.SendErrorResponse(c, http.StatusNotFound, consts.ErrProductNotFound)
		} else {
			helper_handler.SendErrorResponse(c, http.StatusInternalServerError, consts.ErrFailedUpdate)
		}
		return
	}

	productResponse := transformer.TransformProductEntityToResponse(product)
	helper.LogSuccess("Product updated", product.ID(), product.Name()) // Use global helper
	c.JSON(http.StatusOK, productResponse)
}

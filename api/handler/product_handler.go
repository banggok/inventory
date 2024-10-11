package handler

import (
	consts "inventory_management/api/handler/const"
	"inventory_management/api/handler/dto"
	helper_handler "inventory_management/api/handler/helper"
	"inventory_management/api/handler/transformer"
	"inventory_management/internal/usecase"
	"inventory_management/pkg/utility"
	"math"
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
	if validationErrors, err := helper_handler.ReadAndValidateRequestBody(c, &req); validationErrors != nil {
		helper_handler.SendErrorResponse(c, http.StatusUnprocessableEntity, validationErrors)
		return
	} else if err != nil {
		helper_handler.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	product, err := h.productUsecase.CreateProduct(req.Name)
	if err != nil {
		utility.LogError(consts.ErrFailedCreate, req.Name, err)
		helper_handler.SendErrorResponse(c, http.StatusInternalServerError, consts.ErrFailedCreate)
		return
	}

	productResponse := transformer.TransformProductEntityToResponse(product)
	utility.LogSuccess("Product created successfully", product.ID(), product.Name())
	c.JSON(http.StatusCreated, productResponse)
}

// GetProduct retrieves a product by its ID
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || id == 0 {
		utility.LogError(consts.ErrInvalidProductID, idParam, err)
		helper_handler.SendErrorResponse(c, http.StatusBadRequest, consts.ErrInvalidProductID)
		return
	}

	product, err := h.productUsecase.GetProductByID(uint(id))
	if err != nil {
		if err == usecase.ErrProductNotFound { // Use a defined error from the usecase layer
			helper_handler.SendErrorResponse(c, http.StatusNotFound, consts.ErrProductNotFound)
		} else {
			// Ensure the ID fits within int range before converting
			if id <= uint64(math.MaxInt) {
				utility.LogError(consts.ErrFailedRetrieve, strconv.Itoa(int(id)), err)
			} else {
				utility.LogError(consts.ErrFailedRetrieve, idParam, err) // Use original string if too large
			}
			helper_handler.SendErrorResponse(c, http.StatusInternalServerError, consts.ErrFailedRetrieve)
		}
		return
	}

	productResponse := transformer.TransformProductEntityToResponse(product)
	utility.LogSuccess("Product retrieved", product.ID(), product.Name())
	c.JSON(http.StatusOK, productResponse)
}

// UpdateProductName handles updating a product's name
func (h *ProductHandler) UpdateProductName(c *gin.Context) {
	var req dto.UpdateProductRequest
	if validationErrors, err := helper_handler.ReadAndValidateRequestBody(c, &req); validationErrors != nil {
		helper_handler.SendErrorResponse(c, http.StatusUnprocessableEntity, validationErrors)
		return
	} else if err != nil {
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
		if err == usecase.ErrProductNotFound { // Use a defined error from the usecase layer
			helper_handler.SendErrorResponse(c, http.StatusNotFound, consts.ErrProductNotFound)
		} else {
			// Ensure the ID fits within int range before converting
			if id <= uint64(math.MaxInt) {
				utility.LogError(consts.ErrFailedUpdate, strconv.Itoa(int(id)), err)
			} else {
				utility.LogError(consts.ErrFailedUpdate, idParam, err) // Use original string if too large
			}
			helper_handler.SendErrorResponse(c, http.StatusInternalServerError, consts.ErrFailedUpdate)
		}
		return
	}

	productResponse := transformer.TransformProductEntityToResponse(product)
	utility.LogSuccess("Product updated successfully", product.ID(), product.Name())
	c.JSON(http.StatusOK, productResponse)
}

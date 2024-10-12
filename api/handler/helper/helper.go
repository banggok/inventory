package helper_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	consts "inventory_management/api/handler/const"
	"inventory_management/api/handler/dto"
	"inventory_management/pkg/utility"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ReadAndValidateRequestBody reads the request body, validates it, and returns validation errors if any.
// It returns nil if validation succeeds.
func ReadAndValidateRequestBody(c *gin.Context, request dto.Validator) (map[string]string, error) {
	// Read the request body using io.ReadAll instead of ioutil.ReadAll
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, errors.New("invalid request body")
	}

	// Unmarshal the raw body into the request structure
	if err := json.Unmarshal(rawBody, request); err != nil {
		return nil, errors.New("invalid JSON format")
	}

	// Validate the request body after it has been populated
	validationErrors := request.Validate()
	if validationErrors != nil {
		// Return validation errors
		return validationErrors, nil
	}

	// Return nil if no validation errors
	return nil, nil
}

// ParseIDFromParam extracts and validates the ID from URL parameters.
func ParseIDFromParam(c *gin.Context) (uint, error) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || id == 0 {
		return 0, fmt.Errorf("%s: %v", consts.ErrInvalidProductID, err)
	}
	return uint(id), nil
}

// HandleErrorResponse is a reusable function to handle error responses and logging.
func HandleErrorResponse(c *gin.Context, err error, errorMessage string, statusCode int) {
	utility.LogError(errorMessage, "", err)
	c.JSON(statusCode, gin.H{"errors": errorMessage})
}

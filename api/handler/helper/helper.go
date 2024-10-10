package helper_handler

import (
	"encoding/json"
	"inventory_management/api/handler/dto"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

// SendErrorResponse sends a consistent error response. It can handle both string and map errors.
func SendErrorResponse(c *gin.Context, statusCode int, message interface{}) {
	c.JSON(statusCode, gin.H{"errors": message})
}

// ReadAndValidateRequestBody reads the request body, validates it, and returns validation errors if any.
// It returns nil if validation succeeds.
func ReadAndValidateRequestBody(c *gin.Context, request dto.Validator) (map[string]string, error) {
	// Read the request body
	rawBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return map[string]string{"error": "Invalid request body"}, err
	}

	// Unmarshal the raw body into the request structure
	if err := json.Unmarshal(rawBody, request); err != nil {
		return map[string]string{"error": "Invalid JSON format"}, err
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

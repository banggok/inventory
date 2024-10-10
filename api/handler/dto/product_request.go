package dto

import (
	"github.com/go-playground/validator/v10"
)

// CreateProductRequest represents the request body for creating a product
type CreateProductRequest struct {
	Name string `json:"name" validate:"required,min=2,max=255"`
}

// Validate performs JSON decoding and validation on CreateProductRequest and returns custom error messages if validation fails.
func (r *CreateProductRequest) Validate() map[string]string {

	// Create a new validator instance
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return r.parseValidationErrors(err.(validator.ValidationErrors))
	}

	return nil
}

// parseValidationErrors converts the validation errors into a map of custom error messages.
func (r *CreateProductRequest) parseValidationErrors(validationErrors validator.ValidationErrors) map[string]string {
	errors := make(map[string]string)

	for _, err := range validationErrors {
		fieldWithTag := err.Field() + "." + err.Tag()
		errors[err.Field()] = r.getCustomErrorMessage(fieldWithTag)
	}

	return errors
}

// getCustomErrorMessage returns custom error messages for validation rules.
func (r *CreateProductRequest) getCustomErrorMessage(fieldWithTag string) string {
	customMessages := map[string]string{
		"Name.required": "Product name is required.",
		"Name.min":      "Product name must be at least 2 characters long.",
		"Name.max":      "Product name must be less than 255 characters long.",
	}

	if message, exists := customMessages[fieldWithTag]; exists {
		return message
	}
	return "Invalid field"
}

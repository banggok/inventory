package dto

import (
	"net/url"
	"strconv"

	"github.com/go-playground/validator/v10"
)

// ProductListQueryParams defines the query parameters for listing products
type ProductListQueryParams struct {
	SearchTerm    string `json:"search"`
	SortBy        string `json:"sortBy" validate:"oneof=name sku"`
	SortDirection string `json:"sortDirection" validate:"oneof=asc desc"`
	Limit         int    `json:"limit"`
	Offset        int    `json:"offset"`
}

// Validate performs validation on the query parameters and returns custom error messages
func (p *ProductListQueryParams) Validate(queryParams url.Values) map[string]string {
	// Manually extract and assign the query parameters
	p.SearchTerm = queryParams.Get("search")
	p.SortBy = queryParams.Get("sortBy")
	p.SortDirection = queryParams.Get("sortDirection")

	// If limit or offset are not provided, set default values
	if limit := queryParams.Get("limit"); limit != "" {
		p.Limit, _ = strconv.Atoi(limit)
	} else {
		p.Limit = 10 // default value
	}

	if offset := queryParams.Get("offset"); offset != "" {
		p.Offset, _ = strconv.Atoi(offset)
	} else {
		p.Offset = 0 // default value
	}

	// Perform validation using the validator package
	validate := validator.New()
	err := validate.Struct(p)
	if err != nil {
		return p.parseValidationErrors(err.(validator.ValidationErrors))
	}

	return nil
}

// parseValidationErrors converts validation errors into custom error messages
func (p *ProductListQueryParams) parseValidationErrors(validationErrors validator.ValidationErrors) map[string]string {
	errors := make(map[string]string)

	for _, err := range validationErrors {
		fieldWithTag := err.Field() + "." + err.Tag()
		errors[err.Field()] = p.getCustomErrorMessage(fieldWithTag)
	}

	return errors
}

// getCustomErrorMessage returns custom error messages based on the field and tag
func (p *ProductListQueryParams) getCustomErrorMessage(fieldWithTag string) string {
	customMessages := map[string]string{
		"SortBy.oneof":        "sortBy must be either 'name' or 'sku'.",
		"SortDirection.oneof": "sortDirection must be either 'asc' or 'desc'.",
	}

	if message, exists := customMessages[fieldWithTag]; exists {
		return message
	}

	return "Invalid field"
}

// /api/handler/dto/product_request.go
package dto

import "errors"

// CreateProductRequest represents the request body for creating a product
type CreateProductRequest struct {
	Name string `json:"name" binding:"required"` // Add binding tags for validation
}

// Validate checks the validity of CreateProductRequest fields
func (req *CreateProductRequest) Validate() error {
	if req.Name == "" {
		return errors.New("name cannot be empty")
	}
	return nil
}

// /api/handler/dto/product_request.go
package dto

// CreateProductRequest represents the request body for creating a product
type CreateProductRequest struct {
	Name string `json:"name" binding:"required"` // Add binding tags for validation
}

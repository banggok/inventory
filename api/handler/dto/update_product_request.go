package dto

// UpdateProductRequest represents the request body for updating the product name.
type UpdateProductRequest struct {
	Name string `json:"name" binding:"required,min=2,max=255"` // Add length validation if needed
}

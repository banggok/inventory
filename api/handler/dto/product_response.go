// /api/handler/dto/product_response.go
package dto

// ProductResponse represents the response body for a product
type ProductResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	SKU  string `json:"sku"`
}

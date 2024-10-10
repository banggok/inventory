package dto

import "time"

// ProductResponse represents the response body for a product
type ProductResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	SKU       string    `json:"sku"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

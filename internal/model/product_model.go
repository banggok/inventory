// /internal/model/product_model.go
package model

import "time"

// Product represents the structure of the products table in the database
type Product struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"` // Matches SERIAL PRIMARY KEY in SQL
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	SKU       string    `gorm:"type:varchar(100);unique;not null" json:"sku"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // Automatically handle creation timestamp
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"` // Automatically handle update timestamp
}

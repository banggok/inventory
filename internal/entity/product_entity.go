// /internal/entity/product.go
package entity

import (
	"errors"
	"math/rand"
	"strconv"
	"time"
)

// Product represents the business logic of a product
type Product struct {
	ID   uint
	Name string
	SKU  string
}

// BeforeCreate is a GORM hook to generate SKU before creating a product
func (p *Product) BeforeCreate() error {
	if err := p.Validate(); err != nil {
		return err
	}
	p.SKU = generateSKU()
	return nil
}

// Validate checks if the product fields are valid
func (p *Product) Validate() error {
	if p.Name == "" {
		return errors.New("name cannot be empty")
	}
	return nil
}

// generateSKU generates a random SKU
func generateSKU() string {
	rand.Seed(time.Now().UnixNano())
	return "SKU-" + strconv.Itoa(rand.Intn(100000)) // Example SKU generation logic
}

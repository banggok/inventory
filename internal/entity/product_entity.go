// /internal/entity/product_entity.go
package entity

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Product represents the business logic of a product
type Product struct {
	id   uint   // Unexported ID field
	name string // Unexported Name field
	sku  string // Unexported SKU field
}

// NewProduct creates a new Product instance and initializes the Name and SKU
func NewProduct(name string) (*Product, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty") // Directly check for empty name
	}

	product := &Product{
		id:   0, // Assign default ID (can change as needed)
		name: name,
		sku:  generateSKU(name), // Generate SKU when creating the product
	}

	return product, nil
}

// MakeProduct sets all attributes of the Product from parameters
func (p *Product) MakeProduct(id uint, name string, sku string) error {
	if name == "" {
		return errors.New("name cannot be empty") // Check for empty name
	}
	p.id = id     // Set the unexported ID
	p.name = name // Set the unexported Name
	p.sku = sku   // Set the SKU from the parameter
	return nil
}

// generateSKU generates an SKU based on the product name and a random number
func generateSKU(name string) string {
	// Use rand.New with a source based on time to avoid deprecation
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Get the first 3 letters of the name, convert to uppercase
	namePart := strings.ToUpper(name)
	if len(namePart) > 3 {
		namePart = namePart[:3] // Take the first 3 characters
	}

	// Generate a random number between 10000 and 99999 for the SKU
	randomPart := r.Intn(90000) + 10000 // Ensures it's a 5-digit number

	// Combine the name part and the random part to form the SKU
	return "SKU-" + namePart + "-" + strconv.Itoa(randomPart)
}

// ID returns the ID of the product
func (p *Product) ID() uint {
	return p.id // Getter for ID
}

// Name returns the Name of the product
func (p *Product) Name() string {
	return p.name // Getter for Name
}

// SKU returns the SKU of the product
func (p *Product) SKU() string {
	return p.sku // Getter for SKU
}

// SetName sets the Name of the product
func (p *Product) SetName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty") // Directly check for empty name
	}
	p.name = name // Set the unexported Name
	return nil
}

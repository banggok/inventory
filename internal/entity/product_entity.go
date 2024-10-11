package entity

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Define constant for error message
const ErrEmptyName = "name cannot be empty"

// Product represents the business logic of a product
type Product struct {
	id        uint      // Unexported ID field
	name      string    // Unexported Name field
	sku       string    // Unexported SKU field
	createdAt time.Time // Unexported CreatedAt field
	updatedAt time.Time // Unexported UpdatedAt field
}

// NewProduct creates a new Product instance and initializes the Name, SKU, and timestamps
func NewProduct(name string) (*Product, error) {
	return NewProductWithCustomGenerator(name, rand.Read) // Default to using rand.Read for random number generation
}

// generateSKUWithCustomGenerator generates an SKU with a custom random function
func generateSKUWithCustomGenerator(name string, randomNumberGenerator func([]byte) (int, error)) (string, error) {
	// Get the first 3 letters of the name, convert to uppercase
	namePart := strings.ToUpper(name)
	if len(namePart) > 3 {
		namePart = namePart[:3] // Take the first 3 characters
	}

	// Generate a secure random 5-digit number
	randomPart, err := generateSecureFiveDigitNumber(randomNumberGenerator)
	if err != nil {
		return "", err
	}

	// Combine the name part and the random part to form the SKU
	return "SKU-" + namePart + "-" + randomPart, nil
}

// NewProductWithCustomGenerator creates a new Product with a custom random number generator (for testing)
func NewProductWithCustomGenerator(name string, randomNumberGenerator func([]byte) (int, error)) (*Product, error) {
	if name == "" {
		return nil, errors.New(ErrEmptyName) // Use constant for empty name check
	}

	currentTime := time.Now()

	sku, err := generateSKUWithCustomGenerator(name, randomNumberGenerator)
	if err != nil {
		return nil, err
	}
	product := &Product{
		id:        0, // Assign default ID (can change as needed)
		name:      name,
		sku:       sku, // Generate SKU when creating the product
		createdAt: currentTime,
		updatedAt: currentTime,
	}

	return product, nil
}

// MakeProduct sets all attributes of the Product from parameters
func (p *Product) MakeProduct(id uint, name string, sku string, createdAt, updatedAt time.Time) error {
	if name == "" {
		return errors.New(ErrEmptyName) // Use constant for empty name check
	}
	p.id = id               // Set the unexported ID
	p.name = name           // Set the unexported Name
	p.sku = sku             // Set the SKU from the parameter
	p.createdAt = createdAt // Set the CreatedAt field
	p.updatedAt = updatedAt // Set the UpdatedAt field
	return nil
}

// generateSecureFiveDigitNumber generates a cryptographically secure 5-digit number
func generateSecureFiveDigitNumber(randomNumberGenerator func([]byte) (int, error)) (string, error) {
	// Create a byte slice for random bytes
	b := make([]byte, 8) // 8 bytes for a 64-bit unsigned int

	// Read cryptographically secure random bytes
	_, err := randomNumberGenerator(b)
	if err != nil {
		return "", err
	}

	// Convert the bytes to a large integer
	randomInt := binary.BigEndian.Uint64(b)

	// Take the last 5 digits of the random integer
	randomString := fmt.Sprintf("%05d", randomInt%100000)

	return randomString, nil
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

// CreatedAt returns the creation timestamp of the product
func (p *Product) CreatedAt() time.Time {
	return p.createdAt
}

// UpdatedAt returns the last updated timestamp of the product
func (p *Product) UpdatedAt() time.Time {
	return p.updatedAt
}

// SetName sets the Name of the product
func (p *Product) SetName(name string) error {
	if name == "" {
		return errors.New(ErrEmptyName) // Use constant for empty name check
	}
	p.name = name // Set the unexported Name
	return nil
}

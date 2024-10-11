package entity_test

import (
	"inventory_management/internal/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNewProduct tests the NewProduct function
// TestNewProduct tests the NewProduct function
func TestNewProduct(t *testing.T) {
	// Test valid product creation
	product, err := entity.NewProduct("TestProduct")
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, "TestProduct", product.Name())

	// Ensure the SKU follows the expected format (e.g., SKU-XXX-XXXXX)
	assert.Contains(t, product.SKU(), "SKU-")
	assert.Len(t, product.SKU(), 13) // Ensure SKU is in the expected format (9 chars + 5 digits)

	// Ensure timestamps are set correctly
	assert.WithinDuration(t, time.Now(), product.CreatedAt(), time.Second)
	assert.WithinDuration(t, time.Now(), product.UpdatedAt(), time.Second)

	// Test error when name is empty
	product, err = entity.NewProduct("")
	assert.Error(t, err)
	assert.Nil(t, product)
	assert.EqualError(t, err, entity.ErrEmptyName)
}

// TestMakeProduct tests the MakeProduct method
func TestMakeProduct(t *testing.T) {
	product := &entity.Product{}
	createdAt := time.Now().Add(-24 * time.Hour)
	updatedAt := time.Now()

	// Test successful MakeProduct
	err := product.MakeProduct(1, "UpdatedProduct", "SKU-UPC-12345", createdAt, updatedAt)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), product.ID())
	assert.Equal(t, "UpdatedProduct", product.Name())
	assert.Equal(t, "SKU-UPC-12345", product.SKU())
	assert.Equal(t, createdAt, product.CreatedAt())
	assert.Equal(t, updatedAt, product.UpdatedAt())

	// Test error when name is empty in MakeProduct
	err = product.MakeProduct(1, "", "SKU-UPC-12345", createdAt, updatedAt)
	assert.Error(t, err)
	assert.EqualError(t, err, entity.ErrEmptyName)
}

// TestGetters tests all the getters (ID, Name, SKU, CreatedAt, UpdatedAt)
func TestGetters(t *testing.T) {
	createdAt := time.Now().Add(-24 * time.Hour)
	updatedAt := time.Now()
	product := &entity.Product{}
	err := product.MakeProduct(1, "TestProduct", "SKU-TP-54321", createdAt, updatedAt)
	assert.NoError(t, err)

	assert.Equal(t, uint(1), product.ID())
	assert.Equal(t, "TestProduct", product.Name())
	assert.Equal(t, "SKU-TP-54321", product.SKU())
	assert.Equal(t, createdAt, product.CreatedAt())
	assert.Equal(t, updatedAt, product.UpdatedAt())
}

// TestSetName tests the SetName method
func TestSetName(t *testing.T) {
	product := &entity.Product{}

	// Test setting a valid name
	err := product.SetName("NewName")
	assert.NoError(t, err)
	assert.Equal(t, "NewName", product.Name())

	// Test error when setting an empty name
	err = product.SetName("")
	assert.Error(t, err)
	assert.EqualError(t, err, entity.ErrEmptyName)
}

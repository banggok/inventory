package repository_test

import (
	"errors"
	"inventory_management/internal/entity"
	"inventory_management/internal/model"
	"inventory_management/internal/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDB is a mock for the DB interface in the repository package
type MockDB struct {
	mock.Mock
}

// Mock Save function for gorm.DB
func (m *MockDB) Save(value interface{}) *gorm.DB {
	args := m.Called(value)
	return &gorm.DB{Error: args.Error(0)}
}

// Mock First function for gorm.DB
func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(append([]interface{}{dest}, conds...)...)
	return &gorm.DB{Error: args.Error(0)}
}

// TestPostgresProductRepository_Save tests the Save method
func TestPostgresProductRepository_Save(t *testing.T) {
	mockDB := new(MockDB) // Fresh mock object for this test
	repo := repository.NewPostgresProductRepository(mockDB)

	// Mock product entity to save
	product := &entity.Product{}
	err := product.MakeProduct(0, "Test Product", "SKU-TST-12345", time.Now(), time.Now())
	assert.NoError(t, err)

	// Simulate successful save
	mockDB.On("Save", mock.Anything).Return(nil)

	// Call the repository save method
	err = repo.Save(product)
	assert.NoError(t, err)

	// Ensure the mock expectations were met
	mockDB.AssertExpectations(t)

	// Test DB save error
	mockDB = new(MockDB) // Re-initialize for a clean state
	repo = repository.NewPostgresProductRepository(mockDB)

	mockDB.On("Save", mock.Anything).Return(errors.New("save error"))
	err = repo.Save(product)
	assert.Error(t, err)
	assert.EqualError(t, err, "save error")

	// Ensure the mock expectations were met
	mockDB.AssertExpectations(t)
}

// TestPostgresProductRepository_FindByID tests the FindByID method
func TestPostgresProductRepository_FindByID(t *testing.T) {
	mockDB := new(MockDB) // Fresh mock object for this test
	repo := repository.NewPostgresProductRepository(mockDB)

	// Create a mock model product to be returned from the database
	modelProduct := &model.Product{
		ID:        1,
		Name:      "Test Product",
		SKU:       "SKU-TST-12345",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test successful retrieval
	mockDB.On("First", mock.Anything, uint(1)).Return(nil).Run(func(args mock.Arguments) {
		dest := args.Get(0).(*model.Product)
		*dest = *modelProduct
	})
	product, err := repo.FindByID(1)
	assert.NoError(t, err)
	assert.Equal(t, modelProduct.Name, product.Name())

	// Ensure the mock expectations were met
	mockDB.AssertExpectations(t)

	// Test "product not found" error
	mockDB = new(MockDB) // Re-initialize for a clean state
	repo = repository.NewPostgresProductRepository(mockDB)

	mockDB.On("First", mock.Anything, uint(999)).Return(gorm.ErrRecordNotFound)
	product, err = repo.FindByID(999)
	assert.Nil(t, product)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrProductNotFound, err)

	// Ensure the mock expectations were met
	mockDB.AssertExpectations(t)

	// Test other DB errors
	mockDB = new(MockDB) // Re-initialize for a clean state
	repo = repository.NewPostgresProductRepository(mockDB)

	mockDB.On("First", mock.Anything, uint(2)).Return(errors.New("db error"))
	product, err = repo.FindByID(2)
	assert.Nil(t, product)
	assert.Error(t, err)
	assert.EqualError(t, err, "db error")

	// Ensure the mock expectations were met
	mockDB.AssertExpectations(t)
}

// TestPostgresProductRepository_Save tests the Save method including MakeProduct error case
func TestPostgresProductRepository_Save_MakeProductError(t *testing.T) {
	mockDB := new(MockDB)
	repo := repository.NewPostgresProductRepository(mockDB)

	// Mock product entity to save
	product := &entity.Product{}
	err := product.MakeProduct(0, "Test Product", "SKU-TST-12345", time.Now(), time.Now())
	assert.NoError(t, err)

	// Simulate successful save
	mockDB.On("Save", mock.Anything).Return(nil)

	// Modify the product to simulate MakeProduct returning an error (e.g., empty name)
	product = &entity.Product{} // Re-create a product with an empty name
	err = product.MakeProduct(0, "", "SKU-TST-12345", time.Now(), time.Now())
	assert.Error(t, err)

	// Call the repository save method and expect an error from MakeProduct
	err = repo.Save(product)
	assert.Error(t, err)
	assert.EqualError(t, err, entity.ErrEmptyName)

	// Ensure the mock expectations were met
	mockDB.AssertExpectations(t)
}

// TestPostgresProductRepository_FindByID tests the FindByID method including modelToEntity error case
func TestPostgresProductRepository_FindByID_ModelToEntityError(t *testing.T) {
	mockDB := new(MockDB)
	repo := repository.NewPostgresProductRepository(mockDB)

	// Mock model product to return from DB
	modelProduct := &model.Product{
		ID:        1,
		Name:      "", // Empty name to trigger MakeProduct error
		SKU:       "SKU-TST-12345",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test case where MakeProduct fails due to an empty name in modelToEntity
	mockDB.On("First", mock.Anything, uint(1)).Return(nil).Run(func(args mock.Arguments) {
		dest := args.Get(0).(*model.Product)
		*dest = *modelProduct
	})
	product, err := repo.FindByID(1)
	assert.Error(t, err)
	assert.Nil(t, product)
	assert.EqualError(t, err, entity.ErrEmptyName)

	// Ensure the mock expectations were met
	mockDB.AssertExpectations(t)
}

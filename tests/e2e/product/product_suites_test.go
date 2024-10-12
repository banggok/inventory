package product_e2e_test

import (
	"inventory_management/internal/entity"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestProductE2E(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "E2E  Product Handler Suite")
}

// Helper function to truncate tables between tests
func TruncateTables(database *gorm.DB) {
	database.Exec("TRUNCATE TABLE products RESTART IDENTITY CASCADE;")
}

// MockProductUsecase is the mock implementation of the ProductUsecase interface.
type MockProductUsecase struct {
	mock.Mock
}

// ListProducts mock method
func (m *MockProductUsecase) ListProducts(searchTerm string, sortBy string, sortDirection string, limit int, offset int) ([]*entity.Product, error) {
	args := m.Called(searchTerm, sortBy, sortDirection, limit, offset)
	if args.Get(0) != nil {
		return args.Get(0).([]*entity.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

// CreateProduct mock method
func (m *MockProductUsecase) CreateProduct(name string) (*entity.Product, error) {
	args := m.Called(name)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

// GetProductByID mock method (missing in your previous implementation)
func (m *MockProductUsecase) GetProductByID(id uint) (*entity.Product, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

// UpdateProductName mock method
func (m *MockProductUsecase) UpdateProductName(id uint, name string) (*entity.Product, error) {
	args := m.Called(id, name)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

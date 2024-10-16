// /internal/usecase/product_usecase.go
package usecase

import (
	"inventory_management/internal/entity"
	"inventory_management/internal/repository"
)

type ProductUsecase interface {
	CreateProduct(name string) (*entity.Product, error)
	GetProductByID(id uint) (*entity.Product, error)
	UpdateProductName(id uint, name string) (*entity.Product, error)
	ListProducts(searchTerm string, sortBy string, sortDirection string, limit int, offset int) ([]*entity.Product, error)
}

type productUsecase struct {
	productRepo repository.PostgresProductRepository
}

func NewProductUsecase(repo repository.PostgresProductRepository) ProductUsecase {
	return &productUsecase{productRepo: repo}
}

func (u *productUsecase) CreateProduct(name string) (*entity.Product, error) {
	p, err := entity.NewProduct(name)
	if err != nil {
		return nil, err
	}
	err = u.productRepo.Save(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (u *productUsecase) GetProductByID(id uint) (*entity.Product, error) {
	product, err := u.productRepo.FindByID(id)
	if err != nil {
		if err == repository.ErrProductNotFound {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return product, nil
}

// UpdateProductName updates the name of an existing product
func (u *productUsecase) UpdateProductName(id uint, name string) (*entity.Product, error) {
	product, err := u.productRepo.FindByID(id)
	if err != nil {
		if err == repository.ErrProductNotFound {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	// Update the product name
	if err := product.SetName(name); err != nil {
		return nil, err
	}
	if err := u.productRepo.Save(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (u *productUsecase) ListProducts(searchTerm string, sortBy string, sortDirection string, limit int, offset int) ([]*entity.Product, error) {
	return u.productRepo.ListProducts(searchTerm, sortBy, sortDirection, limit, offset)
}

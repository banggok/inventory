// /internal/usecase/product_usecase.go
package usecase

import (
	"inventory_management/internal/entity"
	"inventory_management/internal/repository"
)

type ProductUsecase interface {
	CreateProduct(name string) (*entity.Product, error)
	GetProductByID(id uint) (*entity.Product, error)
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
	return u.productRepo.FindByID(id)
}

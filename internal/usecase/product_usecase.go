// /internal/usecase/product/usecase.go
package usecase

import (
	"inventory_management/api/handler/dto"
	"inventory_management/internal/entity"
	"inventory_management/internal/repository"
)

type ProductUsecase interface {
	CreateProduct(dto dto.CreateProductRequest) (*entity.Product, error)
	GetProductByID(id uint) (*entity.Product, error)
}

type productUsecase struct {
	productRepo repository.PostgresProductRepository
}

func NewProductUsecase(repo repository.PostgresProductRepository) ProductUsecase {
	return &productUsecase{productRepo: repo}
}

func (u *productUsecase) CreateProduct(dto dto.CreateProductRequest) (*entity.Product, error) {
	p := &entity.Product{
		Name: dto.Name,
	}
	p.BeforeCreate()
	err := u.productRepo.Save(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (u *productUsecase) GetProductByID(id uint) (*entity.Product, error) {
	return u.productRepo.FindByID(id)
}

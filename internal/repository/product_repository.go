// /internal/repository/product/postgres_product_repository.go
package repository

import (
	"inventory_management/internal/entity"

	"gorm.io/gorm"
)

type PostgresProductRepository interface {
	Save(p *entity.Product) error
	FindByID(id uint) (*entity.Product, error)
}

type postgresProductRepository struct {
	DB *gorm.DB
}

func NewPostgresProductRepository(db *gorm.DB) PostgresProductRepository {
	return &postgresProductRepository{DB: db}
}

func (r *postgresProductRepository) Save(p *entity.Product) error {
	return r.DB.Create(p).Error
}

func (r *postgresProductRepository) FindByID(id uint) (*entity.Product, error) {
	var prod entity.Product
	err := r.DB.First(&prod, id).Error
	return &prod, err
}

// /internal/repository/product_repository.go
package repository

import (
	"errors"
	"inventory_management/internal/entity"
	"inventory_management/internal/model"

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

// Save converts entity to model, saves it to the database, and updates the entity with the generated values
func (r *postgresProductRepository) Save(p *entity.Product) error {
	// Convert entity.Product to model.Product for database operations
	modelProduct := entityToModel(p)

	// Use GORM's Save method to save the product to the database
	// This will insert or update the record, and auto-populate fields like ID if it's a new record
	if err := r.DB.Save(modelProduct).Error; err != nil {
		return err
	}

	// Update the passed entity.Product with the saved model's ID and other fields
	p.MakeProduct(modelProduct.ID, modelProduct.Name, modelProduct.SKU)

	return nil
}

// FindByID fetches a product from the database, converts model to entity, and returns it
func (r *postgresProductRepository) FindByID(id uint) (*entity.Product, error) {
	var modelProduct model.Product
	err := r.DB.First(&modelProduct, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return a "product not found" error if the product doesn't exist
			return nil, errors.New("product not found")
		}
		// Return other errors as general database errors
		return nil, err
	}
	return modelToEntity(&modelProduct), nil
}

// Convert entity.Product to model.Product for saving to the database
func entityToModel(entityProduct *entity.Product) *model.Product {
	return &model.Product{
		ID:   entityProduct.ID(),
		Name: entityProduct.Name(),
		SKU:  entityProduct.SKU(),
	}
}

// Convert model.Product to entity.Product for returning from the database
func modelToEntity(modelProduct *model.Product) *entity.Product {
	entityProduct := &entity.Product{}
	entityProduct.MakeProduct(modelProduct.ID, modelProduct.Name, modelProduct.SKU)
	return entityProduct
}

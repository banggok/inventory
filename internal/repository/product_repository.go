package repository

import (
	"errors"
	"inventory_management/internal/entity"
	"inventory_management/internal/model"

	"gorm.io/gorm"
)

// Define an interface for the methods we use from gorm.DB
type DB interface {
	Save(value interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
}

// ErrProductNotFound is returned when a product is not found in the database
var ErrProductNotFound = errors.New("product not found")

type PostgresProductRepository interface {
	Save(p *entity.Product) error
	FindByID(id uint) (*entity.Product, error)
}

type postgresProductRepository struct {
	DB DB // Use the interface instead of the concrete gorm.DB type
}

func NewPostgresProductRepository(db DB) PostgresProductRepository {
	return &postgresProductRepository{DB: db}
}

// Save converts entity to model, saves it to the database, and updates the entity with the generated values
func (r *postgresProductRepository) Save(p *entity.Product) error {
	modelProduct := entityToModel(p)

	if err := r.DB.Save(modelProduct).Error; err != nil {
		return err
	}

	if err := p.MakeProduct(
		modelProduct.ID,
		modelProduct.Name,
		modelProduct.SKU,
		modelProduct.CreatedAt,
		modelProduct.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

// FindByID fetches a product from the database, converts model to entity, and returns it
func (r *postgresProductRepository) FindByID(id uint) (*entity.Product, error) {
	var modelProduct model.Product
	err := r.DB.First(&modelProduct, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return modelToEntity(&modelProduct)
}

// Convert entity.Product to model.Product for saving to the database
func entityToModel(entityProduct *entity.Product) *model.Product {
	return &model.Product{
		ID:        entityProduct.ID(),
		Name:      entityProduct.Name(),
		SKU:       entityProduct.SKU(),
		CreatedAt: entityProduct.CreatedAt(),
		UpdatedAt: entityProduct.UpdatedAt(),
	}
}

// Convert model.Product to entity.Product for returning from the database
func modelToEntity(modelProduct *model.Product) (*entity.Product, error) {
	entityProduct := &entity.Product{}

	if err := entityProduct.MakeProduct(
		modelProduct.ID,
		modelProduct.Name,
		modelProduct.SKU,
		modelProduct.CreatedAt,
		modelProduct.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return entityProduct, nil
}

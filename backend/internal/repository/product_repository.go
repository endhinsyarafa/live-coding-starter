package repository

import (
	"live-coding-starter/backend/internal/model"

	"gorm.io/gorm"
)

type ProductRepository interface {
	FindAll() ([]model.Product, error)
	FindByID(id uint) (*model.Product, error)
	Create(product *model.Product) error
	Update(product *model.Product) error
	Delete(product *model.Product) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindAll() ([]model.Product, error) {
	var products []model.Product
	err := r.db.Order("created_at DESC").Find(&products).Error

	return products, err
}

func (r *productRepository) FindByID(id uint) (*model.Product, error) {
	var product model.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(product *model.Product) error {
	return r.db.Delete(product).Error
}

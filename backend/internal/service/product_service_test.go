package service

import (
	"errors"
	"testing"

	"live-coding-starter/backend/internal/dto"
	"live-coding-starter/backend/internal/model"

	"gorm.io/gorm"
)

type fakeProductRepository struct {
	products map[uint]*model.Product
	nextID   uint
	findErr  error
}

func newFakeProductRepository() *fakeProductRepository {
	return &fakeProductRepository{
		products: map[uint]*model.Product{},
		nextID:   1,
	}
}

func (r *fakeProductRepository) FindAll() ([]model.Product, error) {
	products := make([]model.Product, 0, len(r.products))
	for _, product := range r.products {
		products = append(products, *product)
	}

	return products, nil
}

func (r *fakeProductRepository) FindByID(id uint) (*model.Product, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}

	product, ok := r.products[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}

	copiedProduct := *product
	return &copiedProduct, nil
}

func (r *fakeProductRepository) Create(product *model.Product) error {
	product.ID = r.nextID
	r.nextID++

	copiedProduct := *product
	r.products[product.ID] = &copiedProduct

	return nil
}

func (r *fakeProductRepository) Update(product *model.Product) error {
	copiedProduct := *product
	r.products[product.ID] = &copiedProduct

	return nil
}

func (r *fakeProductRepository) Delete(product *model.Product) error {
	delete(r.products, product.ID)

	return nil
}

func TestCreateProductTrimsName(t *testing.T) {
	repository := newFakeProductRepository()
	service := NewProductService(repository)

	product, err := service.CreateProduct(dto.CreateProductRequest{
		Name:        "  Keyboard Mechanical  ",
		Description: "Mechanical keyboard",
		Price:       750000,
		Stock:       10,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.Name != "Keyboard Mechanical" {
		t.Fatalf("expected trimmed product name, got %q", product.Name)
	}
}

func TestCreateProductRejectsBlankName(t *testing.T) {
	repository := newFakeProductRepository()
	service := NewProductService(repository)

	_, err := service.CreateProduct(dto.CreateProductRequest{
		Name:  "   ",
		Price: 750000,
		Stock: 10,
	})

	if !errors.Is(err, ErrInvalidProduct) {
		t.Fatalf("expected ErrInvalidProduct, got %v", err)
	}
}

func TestGetProductByIDMapsNotFound(t *testing.T) {
	repository := newFakeProductRepository()
	service := NewProductService(repository)

	_, err := service.GetProductByID(99)

	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

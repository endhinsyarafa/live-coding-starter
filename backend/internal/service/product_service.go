package service

import (
	"errors"
	"strings"

	"live-coding-starter/backend/internal/dto"
	"live-coding-starter/backend/internal/model"
	"live-coding-starter/backend/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidProduct  = errors.New("invalid product data")
)

type ProductService interface {
	GetAllProducts() ([]model.Product, error)
	GetProductByID(id uint) (*model.Product, error)
	CreateProduct(req dto.CreateProductRequest) (*model.Product, error)
	UpdateProduct(id uint, req dto.UpdateProductRequest) (*model.Product, error)
	DeleteProduct(id uint) error
}

type productService struct {
	repository repository.ProductRepository
}

func NewProductService(repository repository.ProductRepository) ProductService {
	return &productService{repository: repository}
}

func (s *productService) GetAllProducts() ([]model.Product, error) {
	return s.repository.FindAll()
}

func (s *productService) GetProductByID(id uint) (*model.Product, error) {
	product, err := s.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}

		return nil, err
	}

	return product, nil
}

func (s *productService) CreateProduct(req dto.CreateProductRequest) (*model.Product, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrInvalidProduct
	}

	product := &model.Product{
		Name:        name,
		Description: strings.TrimSpace(req.Description),
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if err := s.repository.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) UpdateProduct(id uint, req dto.UpdateProductRequest) (*model.Product, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrInvalidProduct
	}

	product, err := s.GetProductByID(id)
	if err != nil {
		return nil, err
	}

	product.Name = name
	product.Description = strings.TrimSpace(req.Description)
	product.Price = req.Price
	product.Stock = req.Stock

	if err := s.repository.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) DeleteProduct(id uint) error {
	product, err := s.GetProductByID(id)
	if err != nil {
		return err
	}

	return s.repository.Delete(product)
}

package handler

import (
	"errors"
	"net/http"
	"strconv"

	"live-coding-starter/backend/internal/dto"
	"live-coding-starter/backend/internal/service"
	"live-coding-starter/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("", h.GetProducts)
	router.GET("/:id", h.GetProduct)
	router.POST("", h.CreateProduct)
	router.PUT("/:id", h.UpdateProduct)
	router.DELETE("/:id", h.DeleteProduct)
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	products, err := h.service.GetAllProducts()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve products", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Products retrieved successfully", products)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	product, err := h.service.GetProductByID(id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Product retrieved successfully", product)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	product, err := h.service.CreateProduct(req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Product created successfully", product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	product, err := h.service.UpdateProduct(id, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Product updated successfully", product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	if err := h.service.DeleteProduct(id); err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Product deleted successfully", nil)
}

func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(c, http.StatusBadRequest, "Invalid product ID", nil)
		return 0, false
	}

	return uint(id), true
}

func handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrProductNotFound):
		response.Error(c, http.StatusNotFound, "Product not found", nil)
	case errors.Is(err, service.ErrInvalidProduct):
		response.Error(c, http.StatusBadRequest, "Invalid product data", nil)
	default:
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
	}
}

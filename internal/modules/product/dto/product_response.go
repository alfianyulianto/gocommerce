package dto

import "github.com/alfianyulianto/gocommerce/internal/modules/product/entity"

type ProductResponse struct {
	ID          string  `json:"id"`
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Image       string  `json:"image"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Stock       int16   `json:"stock"`
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
	DeletedAt   *int64  `json:"deleted_at"`
}

func ToProductResponse(product *entity.Product) *ProductResponse {
	return &ProductResponse{
		ID:          product.ID,
		SKU:         product.SKU,
		Name:        product.Name,
		Image:       product.Image,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		DeletedAt:   product.DeletedAt,
	}
}

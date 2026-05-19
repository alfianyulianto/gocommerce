package dto

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,lte=255"`
	SKU         string  `json:"sku" validate:"required,alphanum,lte=100"`
	Image       string  `json:"image" validate:"required,lte=255"`
	Description string  `json:"description" validate:"required"`
	Price       float32 `json:"price" validate:"required,number"`
	Stock       int16   `json:"stock" validate:"required,number"`
}

package dto

type CreateOrderRequest struct {
	UserID string             `json:"user_id" validate:"required,uuid"`
	Note   string             `json:"note" validate:"omitempty"`
	Items  []OrderItemRequest `json:"items" validate:"required,min=1"`
}

type OrderItemRequest struct {
	ProductID string `json:"product_id" validate:"required,uuid"`
	Quantity  int16  `json:"quantity" validate:"required,min=1"`
}

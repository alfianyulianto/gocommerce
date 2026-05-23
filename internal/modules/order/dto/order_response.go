package dto

import (
	"fmt"

	"github.com/alfianyulianto/gocommerce/config"
	"github.com/alfianyulianto/gocommerce/internal/modules/order/entity"
)

type OrderResponse struct {
	ID          string              `json:"id"`
	UserID      string              `json:"user_id"`
	Status      string              `json:"status"`
	TotalAmount float32             `json:"total_amount"`
	Note        string              `json:"note"`
	CreatedAt   int64               `json:"created_at"`
	UpdatedAt   int64               `json:"updated_at"`
	Items       []OrderItemResponse `json:"items"`
}

type OrderItemResponse struct {
	ID           string  `json:"id"`
	OrderID      string  `json:"order_id"`
	ProductID    string  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductImage string  `json:"product_image"`
	Quantity     int16   `json:"quantity"`
	UnitPrice    float32 `json:"unit_price"`
	Subtotal     float32 `json:"subtotal"`
	CreatedAt    int64   `json:"created_at"`
}

func ToOrderResponse(order *entity.Order, cfg *config.Config) *OrderResponse {
	response := &OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		Status:      string(order.Status),
		TotalAmount: order.TotalAmount,
		Note:        order.Note,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}

	for _, item := range order.Items {
		response.Items = append(response.Items, OrderItemResponse{
			ID:           item.ID,
			OrderID:      item.OrderID,
			ProductID:    item.ProductID,
			ProductName:  item.ProductName,
			ProductImage: fmt.Sprintf("%s:%d%s", cfg.App.BaseURL, cfg.App.Port, item.ProductImage),
			Quantity:     item.Quantity,
			UnitPrice:    item.UnitPrice,
			Subtotal:     item.Subtotal,
			CreatedAt:    item.CreatedAt,
		})
	}
	return response
}

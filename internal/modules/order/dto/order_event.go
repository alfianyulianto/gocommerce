package dto

import "github.com/alfianyulianto/gocommerce/internal/modules/order/entity"

type OrderEvent struct {
	ID          string           `json:"id"`
	UserID      string           `json:"user_id"`
	Status      string           `json:"status"`
	TotalAmount float32          `json:"total_amount"`
	Note        string           `json:"note"`
	CreatedAt   int64            `json:"created_at"`
	UpdatedAt   int64            `json:"updated_at"`
	Items       []OrderItemEvent `json:"items"`
}

type OrderItemEvent struct {
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

func OrderToEvent(order *entity.Order) *OrderEvent {
	event := &OrderEvent{
		ID:          order.ID,
		UserID:      order.UserID,
		Status:      string(order.Status),
		TotalAmount: order.TotalAmount,
		Note:        order.Note,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}

	for _, item := range order.Items {
		event.Items = append(event.Items, OrderItemEvent{
			ID:           item.ID,
			OrderID:      item.ID,
			ProductID:    item.ProductID,
			ProductName:  item.ProductName,
			ProductImage: item.ProductImage,
			Quantity:     item.Quantity,
			UnitPrice:    item.UnitPrice,
			Subtotal:     item.Subtotal,
			CreatedAt:    item.CreatedAt,
		})
	}

	return event
}

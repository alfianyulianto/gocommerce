package repository

import (
	"context"

	"github.com/alfianyulianto/gocommerce/internal/modules/order/entity"
	"github.com/alfianyulianto/gocommerce/internal/shared"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
	Update(ctx context.Context, order *entity.Order) error
	Delete(ctx context.Context, order *entity.Order) error
	FindById(ctx context.Context, order *entity.Order, id string) error
}

type orderRepository struct {
	shared.Repository[entity.Order]
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		Repository: shared.Repository[entity.Order]{
			DB: db,
		},
	}
}

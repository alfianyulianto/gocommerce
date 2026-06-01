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
	CreateWithItems(ctx context.Context, tx *gorm.DB, order *entity.Order, items []entity.OrderItem) error
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error
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

func (o *orderRepository) CreateWithItems(ctx context.Context, tx *gorm.DB, order *entity.Order, items []entity.OrderItem) error {
	order.Items = items
	if err := tx.WithContext(ctx).Create(order).Error; err != nil {
		return err
	}
	return nil
}

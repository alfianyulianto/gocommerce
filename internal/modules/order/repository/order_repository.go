package repository

import (
	"context"

	"github.com/alfianyulianto/gocommerce/internal/modules/order/entity"
	"github.com/alfianyulianto/gocommerce/internal/shared"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(ctx context.Context, db *gorm.DB, order *entity.Order) error
	Update(ctx context.Context, db *gorm.DB, order *entity.Order) error
	Delete(ctx context.Context, db *gorm.DB, order *entity.Order) error
	FindById(ctx context.Context, db *gorm.DB, order *entity.Order, id string) error
}

type orderRepository struct {
	shared.Repository[entity.Order]
}

func NewOrderRepository() OrderRepository {
	return &orderRepository{}
}

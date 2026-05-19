package repository

import (
	"context"

	"github.com/alfianyulianto/gocommerce/internal/modules/product/entity"
	"github.com/alfianyulianto/gocommerce/internal/shared"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, db *gorm.DB, product *entity.Product) error
	Update(ctx context.Context, db *gorm.DB, product *entity.Product) error
	Delete(ctx context.Context, db *gorm.DB, product *entity.Product) error
	FindById(ctx context.Context, db *gorm.DB, product *entity.Product, id string) error
}

type productRepository struct {
	shared.Repository[entity.Product]
}

func NewProductRepository() ProductRepository {
	return &productRepository{}
}

package repository

import (
	"context"
	"fmt"

	"github.com/alfianyulianto/gocommerce/internal/modules/product/entity"
	"github.com/alfianyulianto/gocommerce/internal/shared"
	"gorm.io/gorm"
)

type ProductFilter struct {
	MinPrice int32 `json:"min_price"`
	MaxPrice int32 `json:"max_price"`
	shared.PaginationFilter
}

type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, product *entity.Product) error
	FindById(ctx context.Context, product *entity.Product, id string) error
	FindAll(ctx context.Context, filter *ProductFilter) ([]entity.Product, int64, error)
	DeductStock(ctx context.Context, tx *gorm.DB, product *entity.Product, quantity int16) error
	UpdateStock(ctx context.Context, productID string, stock int16) error
}

type productRepository struct {
	shared.Repository[entity.Product]
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		Repository: shared.Repository[entity.Product]{
			DB: db,
		},
	}
}

func (p *productRepository) FindAll(ctx context.Context, filter *ProductFilter) ([]entity.Product, int64, error) {
	var products []entity.Product
	err := p.DB.WithContext(ctx).Scopes(p.Filter(filter)).Order(fmt.Sprintf("%s %s", filter.OrderBy, filter.OrderDirection)).Offset((filter.Page - 1) * filter.PerPage).Limit(filter.PerPage).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	var total int64
	err = p.DB.WithContext(ctx).Model(new(entity.Product)).Scopes(p.Filter(filter)).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (p *productRepository) Filter(filter *ProductFilter) func(db *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if filter.MinPrice > 0 {
			tx = tx.Where("price >= ?", filter.MinPrice)
		}

		if filter.MaxPrice > 0 {
			tx = tx.Where("price <= ?", filter.MaxPrice)
		}

		return tx
	}
}

func (p *productRepository) DeductStock(ctx context.Context, tx *gorm.DB, product *entity.Product, quantity int16) error {
	if err := tx.WithContext(ctx).Take(product, "id = ?", product.ID).Error; err != nil {
		return err
	}

	if product.Stock < quantity {
		return &StockError{ProductName: product.Name, Stock: product.Stock, Quantity: quantity}
	}

	if err := tx.WithContext(ctx).Model(product).Update("stock", gorm.Expr("stock - ?", quantity)).Error; err != nil {
		return err
	}
	return nil
}

type StockError struct {
	ProductName     string
	Stock, Quantity int16
}

func (e *StockError) Error() string {
	return fmt.Sprintf("Insufficient stock for product %s. Available %d, requested %d", e.ProductName, e.Stock, e.Quantity)
}

func (p *productRepository) UpdateStock(ctx context.Context, productID string, stock int16) error {
	return p.DB.WithContext(ctx).Model(new(entity.Product)).Where("id = ?", productID).Update("stock", stock).Error
}

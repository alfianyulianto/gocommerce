package search

import (
	"context"

	"github.com/alfianyulianto/gocommerce/internal/infrastucture/elasticsearch"
	"github.com/alfianyulianto/gocommerce/internal/modules/product/entity"
)

type ESSearch interface {
	Index(ctx context.Context, entity *entity.Product) error
}

type esProductSearch struct {
	Client *elasticsearch.Client
}

type ProductDocument struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	SKU         string  `json:"sku"`
	Image       string  `json:"image"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Stock       int16   `json:"stock"`
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
	DeletedAt   *int64  `json:"deleted_at"`
}

func NewEsProductSearch(client *elasticsearch.Client) ESSearch {
	return &esProductSearch{Client: client}
}

func (e *esProductSearch) Index(ctx context.Context, entity *entity.Product) error {
	document := ProductDocument{
		ID:          entity.ID,
		Name:        entity.Name,
		SKU:         entity.SKU,
		Image:       entity.Image,
		Description: entity.Description,
		Price:       entity.Price,
		Stock:       entity.Stock,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		DeletedAt:   entity.DeletedAt,
	}

	return e.Client.Index(ctx, e.Client.ProductIndex(), document.ID, document)
}

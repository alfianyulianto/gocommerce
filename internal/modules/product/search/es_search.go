package search

import (
	"context"
	"encoding/json"

	"github.com/alfianyulianto/gocommerce/internal/infrastucture/elasticsearch"
	"github.com/alfianyulianto/gocommerce/internal/modules/product/dto"
	"github.com/alfianyulianto/gocommerce/internal/modules/product/entity"
)

type ESSearch interface {
	Index(ctx context.Context, entity *entity.Product) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, request *dto.SearchProductRequest) ([]dto.ProductResponse, int64, error)
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

func (e *esProductSearch) Delete(ctx context.Context, id string) error {
	return e.Client.Delete(ctx, e.Client.ProductIndex(), id)
}

func (e *esProductSearch) Search(ctx context.Context, request *dto.SearchProductRequest) ([]dto.ProductResponse, int64, error) {
	var must []any
	if request.Search != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":    request.Search,
				"type":     "best_fields",
				"fields":   []string{"name^2", "description"},
				"operator": "OR",
			},
		})
	}

	if request.MinPrice > 0 {
		must = append(must, map[string]interface{}{
			"range": map[string]interface{}{
				"price": map[string]interface{}{"gte": request.MinPrice},
			},
		})
	}

	if request.MaxPrice > 0 {
		must = append(must, map[string]interface{}{
			"range": map[string]interface{}{
				"price": map[string]interface{}{"lte": request.MaxPrice},
			},
		})
	}

	var boolQuery map[string]interface{}
	if len(must) > 0 {
		boolQuery = map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		}
	} else {
		boolQuery = map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []any{
					map[string]interface{}{
						"match_all": map[string]interface{}{},
					},
				},
			},
		}
	}

	docs, total, err := e.Client.Search(ctx, e.Client.ProductIndex(), map[string]interface{}{
		"query": boolQuery,
		"size":  request.PerPage,
		"from":  (request.Page - 1) * request.PerPage,
		"sort": []map[string]interface{}{
			map[string]interface{}{
				request.OrderBy: map[string]interface{}{
					"order": request.OrderDirection,
				},
			},
		},
	})

	if err != nil {
		return nil, 0, err
	}

	var results []dto.ProductResponse
	for _, doc := range docs {
		var d ProductDocument
		err := json.Unmarshal(doc, &d)
		if err != nil {
			continue
		}
		results = append(results, dto.ProductResponse{
			ID:          d.ID,
			SKU:         d.SKU,
			Name:        d.Name,
			Image:       d.Image,
			Description: d.Description,
			Price:       d.Price,
			Stock:       d.Stock,
			CreatedAt:   d.CreatedAt,
			UpdatedAt:   d.UpdatedAt,
			DeletedAt:   d.DeletedAt,
		})
	}

	return results, total, nil
}

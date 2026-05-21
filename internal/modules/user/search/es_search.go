package search

import (
	"context"
	"encoding/json"

	"github.com/alfianyulianto/gocommerce/internal/infrastructure/elasticsearch"
	"github.com/alfianyulianto/gocommerce/internal/modules/user/dto"
	"github.com/alfianyulianto/gocommerce/internal/modules/user/entity"
)

type ESSearch interface {
	Index(ctx context.Context, entity *entity.User) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, request *dto.SearchUserRequest) ([]dto.UserResponse, int64, error)
}

type esUserSearch struct {
	Client *elasticsearch.Client
}

type UserDocument struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Phone     *string `json:"phone"`
	CreatedAt int64   `json:"created_at"`
	UpdatedAt int64   `json:"updated_at"`
	DeletedAt *int64  `json:"deleted_at"`
}

func NewEsUserSearch(client *elasticsearch.Client) ESSearch {
	return &esUserSearch{Client: client}
}

func (e *esUserSearch) Index(ctx context.Context, entity *entity.User) error {
	document := UserDocument{
		ID:        entity.ID,
		Name:      entity.Name,
		Email:     entity.Email,
		Password:  entity.Password,
		Phone:     entity.Phone,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
	}

	return e.Client.Index(ctx, e.Client.UserIndex(), document.ID, document)
}

func (e *esUserSearch) Delete(ctx context.Context, id string) error {
	return e.Client.Delete(ctx, e.Client.UserIndex(), id)
}

func (e *esUserSearch) Search(ctx context.Context, request *dto.SearchUserRequest) ([]dto.UserResponse, int64, error) {
	var must []any
	if request.Search != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":    request.Search,
				"type":     "best_fields",
				"fields":   []string{"name^3", "email^2", "phone"},
				"operator": "OR",
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
				"must": []interface{}{
					map[string]interface{}{
						"match_all": map[string]interface{}{},
					},
				},
			},
		}
	}

	docs, total, err := e.Client.Search(ctx, e.Client.UserIndex(), map[string]interface{}{
		"query": boolQuery,
		"size":  request.PerPage,
		"from":  (request.Page - 1) * request.PerPage,
		"sort": []interface{}{
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

	var results []dto.UserResponse
	for _, doc := range docs {
		var d UserDocument
		err = json.Unmarshal(doc, &d)
		if err != nil {
			continue
		}
		results = append(results, dto.UserResponse{
			ID:        d.ID,
			Name:      d.Name,
			Email:     d.Email,
			Phone:     d.Phone,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
			DeletedAt: d.DeletedAt,
		})
	}
	return results, total, err
}

package dto

import (
	"github.com/alfianyulianto/gocommerce/internal/shared"
)

type SearchProductRequest struct {
	Search                  string `json:"q" validate:"omitempty,min=1"`
	MinPrice                int32  `json:"min_price" validate:"omitempty,gte=1"`
	MaxPrice                int32  `json:"max_price" validate:"omitempty,gte=1"`
	shared.PaginationFilter `validate:"omitempty"`
}

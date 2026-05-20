package dto

import (
	"github.com/alfianyulianto/gocommerce/internal/modules/product/repository"
	"github.com/alfianyulianto/gocommerce/internal/shared"
)

type ProductFilterRequest struct {
	MinPrice                int32 `json:"min_price" validate:"omitempty,gte=1"`
	MaxPrice                int32 `json:"max_price" validate:"omitempty,gte=1"`
	shared.PaginationFilter `validate:"omitempty"`
}

func ToProductFilter(request *ProductFilterRequest) *repository.ProductFilter {
	return &repository.ProductFilter{
		MinPrice: request.MinPrice,
		MaxPrice: request.MaxPrice,
		PaginationFilter: shared.PaginationFilter{
			Page:           request.Page,
			PerPage:        request.PerPage,
			OrderBy:        request.OrderBy,
			OrderDirection: request.OrderDirection,
		},
	}
}

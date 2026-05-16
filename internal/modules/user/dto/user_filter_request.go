package dto

import (
	"github.com/alfianyulianto/gocommerce/internal/shared"
)

type PaginationRequest struct {
	Page           int    `json:"page" validate:"required,number,gt=0"`
	PerPage        int    `json:"per_page" validate:"required,number,gt=0,oneof=10 25 50 100"`
	OrderBy        string `json:"order_by" validate:"required,oneof=created_at name"`
	OrderDirection string `json:"order_direction" validate:"required,oneof=asc desc"`
}

func (p *PaginationRequest) ToPaginationFilter() *shared.PaginationFilter {
	return &shared.PaginationFilter{
		Page:           p.Page,
		PerPage:        p.PerPage,
		OrderBy:        p.OrderBy,
		OrderDirection: p.OrderDirection,
	}
}

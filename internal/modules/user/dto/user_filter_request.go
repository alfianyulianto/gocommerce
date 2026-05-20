package dto

import (
	"github.com/alfianyulianto/gocommerce/internal/modules/user/repository"
	"github.com/alfianyulianto/gocommerce/internal/shared"
)

type UserFilterRequest struct {
	shared.PaginationFilter `validate:"omitempty"`
}

func ToUserFilter(request *UserFilterRequest) *repository.UserFilter {
	return &repository.UserFilter{
		PaginationFilter: shared.PaginationFilter{
			Page:           request.Page,
			PerPage:        request.PerPage,
			OrderBy:        request.OrderBy,
			OrderDirection: request.OrderDirection,
		},
	}
}

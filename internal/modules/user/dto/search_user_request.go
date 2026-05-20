package dto

import "github.com/alfianyulianto/gocommerce/internal/shared"

type SearchUserRequest struct {
	Search                  string `json:"q" valid:"omitempty,min=1"`
	shared.PaginationFilter `validate:"omitempty"`
}

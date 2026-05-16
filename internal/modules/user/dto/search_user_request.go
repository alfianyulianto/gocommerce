package dto

type SearchUserRequest struct {
	Search string `json:"search" valid:"omitempty,min=1"`
	PaginationRequest
}

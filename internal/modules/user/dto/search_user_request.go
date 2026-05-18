package dto

type SearchUserRequest struct {
	Search string `json:"q" valid:"omitempty,min=1"`
	PaginationRequest
}

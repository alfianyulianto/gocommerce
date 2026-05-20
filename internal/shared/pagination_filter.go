package shared

type PaginationFilter struct {
	Page           int    `json:"page" validate:"required,number,gt=0"`
	PerPage        int    `json:"per_page" validate:"required,number,gt=0,oneof=10 25 50 100"`
	OrderBy        string `json:"order_by" validate:"required,oneof=created_at name"`
	OrderDirection string `json:"order_direction" validate:"required,oneof=asc desc"`
}

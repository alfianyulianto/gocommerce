package shared

type PaginationFilter struct {
	Page           int    `json:"page"`
	PerPage        int    `json:"per_page"`
	OrderBy        string `json:"order_by"`
	OrderDirection string `json:"order_direction"`
}

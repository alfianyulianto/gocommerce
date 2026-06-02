package dto

type UpdateStockRequest struct {
	ID    string `json:"id" valid:"required,uuid"`
	Stock int16  `json:"stock" validate:"required,number"`
}

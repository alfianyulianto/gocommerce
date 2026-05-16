package dto

type UpdateUserRequest struct {
	ID              string  `json:"id" validate:"required,uuid"`
	Name            string  `json:"name" validate:"required,gte=2,lte=255"`
	Email           string  `json:"email" validate:"required,email"`
	Password        *string `json:"password" validate:"omitempty,gte=6,lte=20"`
	ConfirmPassword string  `json:"confirm_password" validate:"omitempty,required_with=Password,gte=6,lte=20,eqfield=Password"`
	Phone           *string `json:"phone" validate:"omitempty,e164,max=20"`
}

package entity

import "github.com/alfianyulianto/gocommerce/internal/modules/user/entity"

type Status string

const (
	StatusPending    Status = "pending"
	StatusConfirmed  Status = "confirmed"
	StatusProcessing Status = "processing"
	StatusShipped    Status = "shipped"
	StatusDelivered  Status = "delivered"
	StatusCanceled   Status = "canceled"
)

type Order struct {
	ID          string  `gorm:"primaryKey;column:id;size:255;<-:create"`
	UserID      string  `gorm:"column:user_id;size:255;<-:create"`
	Status      Status  `gorm:"column:status;size:255;<-:create"`
	TotalAmount float32 `gorm:"column:total_amount;"`
	Note        string  `gorm:"column:note;size:255;"`
	CreatedAt   int64   `gorm:"column:created_at;autoCreateTime;<-:create"`
	UpdatedAt   int64   `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`

	Items []OrderItem `gorm:"foreignKey:OrderID;references:ID"`

	User entity.User `gorm:"foreignKey:UserID;references:ID"`
}

func (o *Order) TableName() string {
	return "orders"
}

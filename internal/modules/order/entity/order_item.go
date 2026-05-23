package entity

import "github.com/alfianyulianto/gocommerce/internal/modules/product/entity"

type OrderItem struct {
	ID           string  `gorm:"primaryKey;column:id;size:255;<-:create"`
	OrderID      string  `gorm:"column:order_id;size:255;<-:create"`
	ProductID    string  `gorm:"column:product_id;size:255;<-:create"`
	ProductName  string  `gorm:"column:product_name;size:255;<-:create"`
	ProductImage string  `gorm:"column:product_image;size:255;<-:create"`
	Quantity     int16   `grom:"column:quantity;not null"`
	UnitPrice    float32 `grom:"column:unit_price;not null"`
	Subtotal     float32 `grom:"column:subtotal;not null"`
	CreatedAt    int64   `grom:"column:created_at;autoCreateTime;<-:create"`

	Order Order `gorm:"foreignKey:OrderID;references:ID"`

	Product entity.Product `gorm:"foreignKey:ProductID;references:ID"`
}

func (o *OrderItem) TableName() string {
	return "order_items"
}

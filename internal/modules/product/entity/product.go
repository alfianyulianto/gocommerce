package entity

type Product struct {
	ID          string  `gorm:"primaryKey;column=id;size:255;<-:create"`
	Name        string  `grom:"column:name;size=255;not null"`
	SKU         string  `grom:"column:sku;size=100;not null,unique"`
	Image       string  `grom:"column:image;not null"`
	Description string  `grom:"column:description;not null"`
	Price       float32 `grom:"column:price;not null"`
	Stock       int16   `grom:"column:stock;not null"`
	CreatedAt   int64   `grom:"column:created_at;autoCreateTime;<-:create"`
	UpdatedAt   int64   `grom:"column:updated_at;autoUpdateTime;autoUpdateTime;<-"`
	DeletedAt   *int64  `grom:"column:deleted_at;index;<-"`
}

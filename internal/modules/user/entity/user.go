package entity

type User struct {
	ID        string  `gorm:"primaryKey;column=id;size:255;<-:create"`
	Name      string  `grom:"column:name;size=255;not null"`
	Email     string  `gorm:"column=email;size=255;not null,unique"`
	Password  string  `gorm:"column=password;size=255;not null"`
	Phone     *string `gorm:"column=phone;size=20;null"`
	CreatedAt int64   `gorm:"column=created_at;autoCreateTime;<-:create"`
	UpdatedAt int64   `gorm:"column:updated_at;autoUpdateTime;<-"`
	DeletedAt *int64  `gorm:"column:deleted_at;index;<-"`
}

func (u *User) TableName() string {
	return "users"
}

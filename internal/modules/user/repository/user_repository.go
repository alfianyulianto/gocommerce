package repository

import (
	"context"
	"fmt"

	"github.com/alfianyulianto/gocommerce/internal/modules/user/entity"
	"github.com/alfianyulianto/gocommerce/internal/shared"
	"gorm.io/gorm"
)

type UserFilter struct {
	shared.PaginationFilter
}

type UserRepository interface {
	Create(ctx context.Context, db *gorm.DB, entity *entity.User) error
	Update(ctx context.Context, db *gorm.DB, entity *entity.User) error
	Delete(ctx context.Context, db *gorm.DB, entity *entity.User) error
	FindById(ctx context.Context, db *gorm.DB, entity *entity.User, id string) error
	FindAll(ctx context.Context, db *gorm.DB, filter *UserFilter) ([]entity.User, int64, error)
}

type userRepository struct {
	shared.Repository[entity.User]
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (u *userRepository) FindAll(ctx context.Context, db *gorm.DB, filter *UserFilter) ([]entity.User, int64, error) {
	var users []entity.User
	err := db.WithContext(ctx).Order(fmt.Sprintf("%s %s", filter.OrderBy, filter.OrderDirection)).Offset((filter.Page - 1) * filter.PerPage).Limit(filter.PerPage).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	var total int64
	err = db.WithContext(ctx).Model(new(entity.User)).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

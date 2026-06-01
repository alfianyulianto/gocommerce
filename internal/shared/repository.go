package shared

import (
	"context"

	"gorm.io/gorm"
)

type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) Create(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Create(entity).Error
}

func (r *Repository[T]) Update(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Save(entity).Error
}

func (r *Repository[T]) Delete(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Unscoped().Delete(entity).Error
}

func (r *Repository[T]) FindById(ctx context.Context, entity *T, id string) error {
	return r.DB.WithContext(ctx).Take(entity, "id = ?", id).Error
}

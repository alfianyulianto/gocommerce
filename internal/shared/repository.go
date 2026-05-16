package shared

import (
	"context"

	"gorm.io/gorm"
)

type Repository[T any] struct {
}

func (r *Repository[T]) Create(ctx context.Context, db *gorm.DB, entity *T) error {
	return db.WithContext(ctx).Create(entity).Error
}

func (r *Repository[T]) Update(ctx context.Context, db *gorm.DB, entity *T) error {
	return db.WithContext(ctx).Save(entity).Error
}

func (r *Repository[T]) Delete(ctx context.Context, db *gorm.DB, entity *T) error {
	return db.WithContext(ctx).Unscoped().Delete(entity).Error
}

func (r *Repository[T]) FindById(ctx context.Context, db *gorm.DB, entity *T, id string) error {
	return db.WithContext(ctx).Take(entity, "id = ?", id).Error
}

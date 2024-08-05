package repository

import (
	"context"
	"go-ms-template-service/model/domain"

	"gorm.io/gorm"
)

type CinemaRepository interface {
	Save(ctx context.Context, tx *gorm.DB, cinema domain.Cinema) (domain.Cinema, error)
	Update(ctx context.Context, tx *gorm.DB, cinema domain.Cinema, CinemaCode string) domain.Cinema
	Delete(ctx context.Context, tx *gorm.DB, cinema domain.Cinema)
	FindByCode(ctx context.Context, tx *gorm.DB, CinemaCode string) (domain.Cinema, error)
	FindAll(ctx context.Context, tx *gorm.DB) ([]domain.Cinema, error)
}

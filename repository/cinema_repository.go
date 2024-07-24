package repository

import (
	"context"
	"database/sql"
	"x1-cinema/model/domain"
)

type CinemaRepository interface {
	Save(ctx context.Context, tx *sql.Tx, cinema domain.Cinema) domain.Cinema
	Update(ctx context.Context, tx *sql.Tx, cinema domain.Cinema) domain.Cinema
	Delete(ctx context.Context, tx *sql.Tx, cinema domain.Cinema)
	FindByCode(ctx context.Context, tx *sql.Tx, CinemaCode string) (domain.Cinema, error)
	FindAll(ctx context.Context, tx *sql.Tx) []domain.Cinema
}

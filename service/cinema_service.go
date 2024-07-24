package service

import (
	"context"
	"x1-cinema/model/web"
)

type CinemaService interface {
	Create(ctx context.Context, request web.CinemaCreateRequest) web.CinemaResponse
	Update(ctx context.Context, request web.CinemaUpdateRequest) web.CinemaResponse
	Delete(ctx context.Context, CinemaCode string)
	FindByCode(ctx context.Context, CinemaCode string) web.CinemaResponse
	FindAll(ctx context.Context) []web.CinemaResponse
}

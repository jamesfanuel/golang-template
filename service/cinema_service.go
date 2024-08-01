package service

import (
	"context"
	"x1-cinema/model/web"
)

type CinemaService interface {
	Create(ctx context.Context, request web.CinemaCreateRequest) web.CinemaResponseCreate
	Update(ctx context.Context, request web.CinemaUpdateRequest, CinemaCode string) web.CinemaResponseUpdate
	Delete(ctx context.Context, CinemaCode string)
	FindByCode(ctx context.Context, CinemaCode string) web.CinemaResponseFind
	FindAll(ctx context.Context) []web.CinemaResponseFind
}

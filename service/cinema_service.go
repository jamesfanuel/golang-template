package service

import (
	"context"
	"net/http"
	"x1-cinema/model/web"
)

type CinemaService interface {
	Create(ctx context.Context, request web.CinemaCreateRequest, urlRequest *http.Request) web.CinemaCreateResponse
	Update(ctx context.Context, request web.CinemaUpdateRequest, CinemaCode string, urlRequest *http.Request) web.CinemaUpdateResponse
	Delete(ctx context.Context, CinemaCode string, urlRequest *http.Request)
	FindByCode(ctx context.Context, CinemaCode string) web.CinemaResponseFind
	FindAll(ctx context.Context) []web.CinemaResponseFind
}

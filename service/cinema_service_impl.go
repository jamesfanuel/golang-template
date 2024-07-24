package service

import (
	"context"
	"database/sql"
	"x1-cinema/helper"
	"x1-cinema/model/domain"
	"x1-cinema/model/web"
	"x1-cinema/repository"

	"github.com/go-playground/validator"
)

type CinemaServiceImpl struct {
	CinemaRepository repository.CinemaRepository
	DB               *sql.DB
	Validate         *validator.Validate
}

func NewCinemaService(cinemaRepository repository.CinemaRepository, DB *sql.DB, validate *validator.Validate) CinemaService {
	return &CinemaServiceImpl{
		CinemaRepository: cinemaRepository,
		DB:               DB,
		Validate:         validate,
	}
}

func (service *CinemaServiceImpl) Create(ctx context.Context, request web.CinemaCreateRequest) web.CinemaResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	cinema := domain.Cinema{
		CinemaCode: request.CinemaCode,
		CinemaName: request.CinemaName,
	}

	cinema = service.CinemaRepository.Save(ctx, tx, cinema)

	return helper.ToCinemaResponse(cinema)
}

func (service *CinemaServiceImpl) Update(ctx context.Context, request web.CinemaUpdateRequest) web.CinemaResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	cinema, err := service.CinemaRepository.FindByCode(ctx, tx, request.CinemaCode)
	helper.PanicIfError(err)

	cinema.CinemaCode = request.CinemaCode
	cinema.CinemaName = request.CinemaName

	cinema = service.CinemaRepository.Update(ctx, tx, cinema)

	return helper.ToCinemaResponse(cinema)
}

func (service *CinemaServiceImpl) Delete(ctx context.Context, CinemaCode string) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	cinema, err := service.CinemaRepository.FindByCode(ctx, tx, CinemaCode)
	helper.PanicIfError(err)

	service.CinemaRepository.Delete(ctx, tx, cinema)
}

func (service *CinemaServiceImpl) FindByCode(ctx context.Context, CinemaCode string) web.CinemaResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	cinema, err := service.CinemaRepository.FindByCode(ctx, tx, CinemaCode)
	helper.PanicIfError(err)

	return helper.ToCinemaResponse(cinema)
}

func (service *CinemaServiceImpl) FindAll(ctx context.Context) []web.CinemaResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	cinema := service.CinemaRepository.FindAll(ctx, tx)

	var cinemaResponses []web.CinemaResponse
	for _, cinema := range cinema {
		cinemaResponses = append(cinemaResponses, helper.ToCinemaResponse(cinema))
	}

	return cinemaResponses
}

package service

import (
	"context"
	// "database/sql"
	"x1-cinema/exception"
	"x1-cinema/helper"
	"x1-cinema/model/domain"
	"x1-cinema/model/web"
	"x1-cinema/repository"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type CinemaServiceImpl struct {
	CinemaRepository repository.CinemaRepository
	DB               *gorm.DB
	Validate         *validator.Validate
}

func NewCinemaService(cinemaRepository repository.CinemaRepository, DB *gorm.DB, validate *validator.Validate) *CinemaServiceImpl {
	return &CinemaServiceImpl{
		CinemaRepository: cinemaRepository,
		DB:               DB,
		Validate:         validate,
	}
}

func (service *CinemaServiceImpl) Create(ctx context.Context, request web.CinemaCreateRequest) web.CinemaResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx := service.DB.Begin()
	// helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	cinema := domain.Cinema{
		CinemaCode:   request.CinemaCode,
		CinemaName:   request.CinemaName,
		ProvinceCode: request.ProvinceCode,
		CityCode:     request.CityCode,
		RegionCode:   request.RegionCode,
		CinemaLevel:  request.CinemaLevel,
	}

	cinema = service.CinemaRepository.Save(ctx, tx, cinema)

	return helper.ToCinemaResponse(cinema)
}

// func (service *CinemaServiceImpl) Update(ctx context.Context, request web.CinemaUpdateRequest, CinemaCode string) web.CinemaResponse {
// 	err := service.Validate.Struct(request)
// 	helper.PanicIfError(err)

// 	tx, err := service.DB.Begin()
// 	helper.PanicIfError(err)
// 	defer helper.CommitOrRollback(tx)

// 	cinema, err := service.CinemaRepository.FindByCode(ctx, tx, CinemaCode)
// 	helper.PanicIfError(err)

// 	cinema.CinemaName = request.CinemaName

// 	cinema = service.CinemaRepository.Update(ctx, tx, cinema, CinemaCode)

// 	return helper.ToCinemaResponse(cinema)
// }

func (service *CinemaServiceImpl) Delete(ctx context.Context, CinemaCode string) {
	tx := service.DB.Begin()
	// helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	cinema, err := service.CinemaRepository.FindByCode(ctx, tx, CinemaCode)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	service.CinemaRepository.Delete(ctx, tx, cinema)
}

func (service *CinemaServiceImpl) FindByCode(ctx context.Context, CinemaCode string) web.CinemaResponse {
	tx := service.DB.Begin()
	// helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	cinema, err := service.CinemaRepository.FindByCode(ctx, tx, CinemaCode)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return helper.ToCinemaResponse(cinema)
}

func (service *CinemaServiceImpl) FindAll(ctx context.Context) []web.CinemaResponse {
	tx := service.DB.Begin()
	// helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	cinema, err := service.CinemaRepository.FindAll(ctx, tx)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	var cinemaResponses []web.CinemaResponse
	for _, cinema := range cinema {
		cinemaResponses = append(cinemaResponses, helper.ToCinemaResponse(cinema))
	}

	return cinemaResponses
}

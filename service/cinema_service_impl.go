package service

import (
	"context"
	"time"

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

func (service *CinemaServiceImpl) Create(ctx context.Context, request web.CinemaCreateRequest) web.CinemaResponseCreate {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx := service.DB.Begin()
	// helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	cinema := domain.Cinema{
		CinemaCode:      request.CinemaCode,
		CinemaName:      request.CinemaName,
		CinemaOwner:     request.CinemaOwner,
		LocationCode:    request.LocationCode,
		ProvinceCode:    request.ProvinceCode,
		CityCode:        request.CityCode,
		RegionCode:      request.RegionCode,
		CompanyCode:     request.CompanyCode,
		CinemaLevel:     request.CinemaLevel,
		OracleCode:      request.OracleCode,
		IsDataMigration: request.IsDataMigration,
		CloseFlag:       request.CloseFlag,
		CloseStart:      request.CloseStart,
		CloseEnd:        request.CloseEnd,
		OperatorEmail:   request.OperatorEmail,
		CreatedBy:       request.CreatedBy,
		CreatedAt:       time.Now(),
		CreatedHostIp:   request.CreatedHostIp,
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cinema, err = service.CinemaRepository.Save(ctx, tx, cinema)

	if err != nil {
		panic(exception.NewDuplicateKeyError(err.Error()))
	}

	return helper.ToCinemaResponseCreate(cinema)
}

func (service *CinemaServiceImpl) Update(ctx context.Context, request web.CinemaUpdateRequest, CinemaCode string) web.CinemaResponseUpdate {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx := service.DB.Begin()
	// helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	cinema := domain.Cinema{
		CinemaCode:      CinemaCode,
		CinemaOwner:     request.CinemaOwner,
		LocationCode:    request.LocationCode,
		ProvinceCode:    request.ProvinceCode,
		CityCode:        request.CityCode,
		RegionCode:      request.RegionCode,
		CompanyCode:     request.CompanyCode,
		CinemaLevel:     request.CinemaLevel,
		OracleCode:      request.OracleCode,
		IsDataMigration: request.IsDataMigration,
		CloseFlag:       request.CloseFlag,
		CloseStart:      request.CloseStart,
		CloseEnd:        request.CloseEnd,
		OperatorEmail:   request.OperatorEmail,
		UpdatedBy:       request.UpdatedBy,
		UpdatedAt:       time.Now(),
		UpdatedHostIp:   request.UpdatedHostIp,
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cinema = service.CinemaRepository.Update(ctx, tx, cinema, CinemaCode)

	return helper.ToCinemaResponseUpdate(cinema)
}

func (service *CinemaServiceImpl) Delete(ctx context.Context, CinemaCode string) {
	tx := service.DB.Begin()
	// helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cinema, err := service.CinemaRepository.FindByCode(ctx, tx, CinemaCode)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	service.CinemaRepository.Delete(ctx, tx, cinema)
}

func (service *CinemaServiceImpl) FindByCode(ctx context.Context, CinemaCode string) web.CinemaResponseFind {
	tx := service.DB.Begin()
	// helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cinema, err := service.CinemaRepository.FindByCode(ctx, tx, CinemaCode)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return helper.ToCinemaResponse(cinema)
}

func (service *CinemaServiceImpl) FindAll(ctx context.Context) []web.CinemaResponseFind {
	tx := service.DB.Begin()
	// helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cinema, err := service.CinemaRepository.FindAll(ctx, tx)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	var cinemaResponses []web.CinemaResponseFind
	for _, cinema := range cinema {
		cinemaResponses = append(cinemaResponses, helper.ToCinemaResponse(cinema))
	}

	return cinemaResponses
}

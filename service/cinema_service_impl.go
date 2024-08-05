package service

import (
	"context"
	"net/http"
	"strings"
	"time"

	// "database/sql"
	"go-ms-template-service/exception"
	"go-ms-template-service/helper"
	"go-ms-template-service/model/domain"
	"go-ms-template-service/model/web"
	"go-ms-template-service/repository"

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

func (service *CinemaServiceImpl) Create(ctx context.Context, request web.CinemaCreateRequest, urlRequest *http.Request) web.CinemaCreateResponse {
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
		CreatedHostIp:   urlRequest.RemoteAddr,
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cinema, err = service.CinemaRepository.Save(ctx, tx, cinema)

	if err != nil {
		if strings.Split(err.Error(), " ")[1] == "1062" {
			panic(exception.NewDuplicateKeyError(err.Error()))
		} else {
			panic(exception.NewUnknownColumnError(err.Error()))
		}
	}

	return helper.ToCinemaCreateResponse(cinema)
}

func (service *CinemaServiceImpl) Update(ctx context.Context, request web.CinemaUpdateRequest, CinemaCode string, urlRequest *http.Request) web.CinemaUpdateResponse {
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
		UpdatedHostIp:   urlRequest.RemoteAddr,
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cinema = service.CinemaRepository.Update(ctx, tx, cinema, CinemaCode)

	return helper.ToCinemaUpdateResponse(cinema)
}

func (service *CinemaServiceImpl) Delete(ctx context.Context, CinemaCode string, urlRequest *http.Request) {
	tx := service.DB.Begin()
	// helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cinema, err := service.CinemaRepository.FindByCode(ctx, tx, CinemaCode)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	cinema.DeletedHostIp = urlRequest.RemoteAddr

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

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
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CinemaServiceImpl struct {
	CinemaRepository repository.CinemaRepository
	DB               *gorm.DB
	Validate         *validator.Validate
	Log              *logrus.Logger
}

func NewCinemaService(cinemaRepository repository.CinemaRepository, DB *gorm.DB, validate *validator.Validate, logrus *logrus.Logger) *CinemaServiceImpl {
	return &CinemaServiceImpl{
		CinemaRepository: cinemaRepository,
		DB:               DB,
		Validate:         validate,
		Log:              logrus,
	}
}

func (service *CinemaServiceImpl) Create(ctx context.Context, request web.CinemaCreateRequest, urlRequest *http.Request) web.CinemaCreateResponse {
	service.Log.Info("Running Create")

	err := service.Validate.Struct(request)
	helper.PanicIfError(err)
	if err != nil {
		service.Log.Error(err)
	}

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
			service.Log.Error(err)
			panic(exception.NewDuplicateKeyError(err.Error()))
		} else {
			service.Log.Error(err)
			panic(exception.NewUnknownColumnError(err.Error()))
		}
	}

	return helper.ToCinemaCreateResponse(cinema)
}

func (service *CinemaServiceImpl) Update(ctx context.Context, request web.CinemaUpdateRequest, CinemaCode string, urlRequest *http.Request) web.CinemaUpdateResponse {
	service.Log.Info("Running Update")

	err := service.Validate.Struct(request)
	helper.PanicIfError(err)
	if err != nil {
		service.Log.Error(err)
	}

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
	service.Log.Info("Running Delete")

	tx := service.DB.Begin()
	// helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cinema, err := service.CinemaRepository.FindByCode(ctx, tx, CinemaCode)
	if err != nil {
		service.Log.Error(err)
		panic(exception.NewNotFoundError(err.Error()))
	}

	cinema.DeletedHostIp = urlRequest.RemoteAddr

	service.CinemaRepository.Delete(ctx, tx, cinema)
}

func (service *CinemaServiceImpl) FindByCode(ctx context.Context, CinemaCode string) web.CinemaResponseFind {
	service.Log.Info("Running Find by Code")

	tx := service.DB.Begin()
	// helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cinema, err := service.CinemaRepository.FindByCode(ctx, tx, CinemaCode)
	if err != nil {
		service.Log.Error(err)
		panic(exception.NewNotFoundError(err.Error()))
	}

	return helper.ToCinemaResponse(cinema)
}

func (service *CinemaServiceImpl) FindAll(ctx context.Context) []web.CinemaResponseFind {
	service.Log.Info("Running Find All")

	tx := service.DB.Begin()
	// helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cinema, err := service.CinemaRepository.FindAll(ctx, tx)
	if err != nil {
		service.Log.Error(err)
		panic(exception.NewNotFoundError(err.Error()))
	}

	var cinemaResponses []web.CinemaResponseFind
	for _, cinema := range cinema {
		cinemaResponses = append(cinemaResponses, helper.ToCinemaResponse(cinema))
	}

	return cinemaResponses
}

func (service *CinemaServiceImpl) SaveToRedis(ctx context.Context) string {
	service.Log.Info("Running Save Screen to Redis")

	client := redis.NewClient(&redis.Options{
		// Addr: "localhost:6379",
		Addr: "172.16.1.133:6379",
		DB:   0,
	})

	tx := service.DB.Begin()
	// helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	ctx, cancel := context.WithTimeout(ctx, 100*time.Second)
	defer cancel()

	// cinemaCode, _ := service.ScreenRepository.FindUniqueCinemaCode(ctx, tx)

	// var theaterScreenList []web.TheaterScreenList
	// for _, cinemaCode := range cinemaCode {
	// 	var screenList []web.ScreenList
	// 	screens, err := service.ScreenRepository.FindByCinemaCode(ctx, tx, cinemaCode)
	// 	if err != nil {
	// 		service.Log.Error(err)
	// 		panic(exception.NewNotFoundError(err.Error()))
	// 	}

	// 	service.Log.Info("Create Key x1_", cinemaCode)
	// 	for _, screen := range screens {
	// 		screenList = append(screenList, helper.ToScreenList(screen))

	// 		err = client.HSet(ctx, "x1_"+cinemaCode, map[string]interface{}{
	// 			"id":           "x1" + "_" + screen.CinemaCode + "_" + strconv.Itoa(screen.ScreenNo),
	// 			"screen_no":    screen.ScreenNo,
	// 			"theater_code": screen.CinemaCode,
	// 		}).Err()

	// 		if err != nil {
	// 			service.Log.Error(err)
	// 		} else {
	// 			service.Log.Info("Success Insert Data to Redis for ", cinemaCode)
	// 		}
	// 	}

	// 	theaterScreenList = append(theaterScreenList, helper.ToTheaterScreenList(cinemaCode, screenList))
	// }

	client.Close()
	return "ok"

}

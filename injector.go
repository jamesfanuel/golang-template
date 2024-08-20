//go:build wireinject
// +build wireinject

package main

import (
	"go-ms-template-service/app"
	"go-ms-template-service/controller"
	"go-ms-template-service/middleware"
	"go-ms-template-service/repository"
	"go-ms-template-service/service"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/wire"
	"github.com/julienschmidt/httprouter"

	_ "github.com/go-sql-driver/mysql"
)

var cinemaSet = wire.NewSet(
	repository.NewCinemaRepository,
	wire.Bind(new(repository.CinemaRepository), new(*repository.CinemaRepositoryImpl)),
	service.NewCinemaService,
	wire.Bind(new(service.CinemaService), new(*service.CinemaServiceImpl)),
	controller.NewCinemaController,
	wire.Bind(new(controller.CinemaController), new(*controller.CinemaControllerImpl)),
)

func InitializedServer() *http.Server {
	wire.Build(
		app.NewLog,
		app.NewDB,
		validator.New,
		cinemaSet,
		app.NewRouter,
		wire.Bind(new(http.Handler), new(*httprouter.Router)),
		middleware.NewAuthMiddleware,
		NewServer,
	)

	return &http.Server{}
}

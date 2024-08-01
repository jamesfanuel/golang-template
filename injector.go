//go:build wireinject
// +build wireinject

package main

import (
	"net/http"
	"x1-cinema/app"
	"x1-cinema/controller"
	"x1-cinema/middleware"
	"x1-cinema/repository"
	"x1-cinema/service"

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

package main

import (
	"net/http"
	"x1-cinema/app"
	"x1-cinema/controller"
	"x1-cinema/helper"
	"x1-cinema/middleware"
	"x1-cinema/repository"
	"x1-cinema/service"

	"github.com/go-playground/validator"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db := app.NewDB()
	validate := validator.New()

	cinemaRepository := repository.NewCinemaRepository()
	cinemaService := service.NewCinemaService(cinemaRepository, db, validate)
	cinemaController := controller.NewCinemaController(cinemaService)

	router := app.NewRouter(cinemaController)

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: middleware.NewAuthMiddleware(router),
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}

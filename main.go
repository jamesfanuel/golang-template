package main

import (
	"fmt"
	"net/http"
	"x1-cinema/app"
	"x1-cinema/controller"
	"x1-cinema/helper"

	// "x1-cinema/helper"
	"x1-cinema/repository"
	"x1-cinema/service"

	"github.com/go-playground/validator"
	_ "github.com/go-sql-driver/mysql"
)

// @title x1-cinema API
// @version 1.0
// @description API Doc for x1-cinema.
// @BasePath /api/v1
func main() {

	app.NewLog("Info", "Initiate Application")
	app.NewEureka()
	db := app.NewDB()
	validate := validator.New()

	cinemaRepository := repository.NewCinemaRepository()
	cinemaService := service.NewCinemaService(cinemaRepository, db, validate)
	cinemaController := controller.NewCinemaController(cinemaService)

	router := app.NewRouter(cinemaController)

	server := http.Server{
		Addr: ":6010",
		// Handler: middleware.NewAuthMiddleware(router),
		Handler: router,
	}

	fmt.Print("Ready To Serve\n")

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}

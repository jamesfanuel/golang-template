package app

import (
	"x1-cinema/controller"
	"x1-cinema/exception"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(cinemaController controller.CinemaController) *httprouter.Router {
	router := httprouter.New()

	router.GET("/api/cinema", cinemaController.FindAll)
	router.GET("/api/cinema/:companyCode", cinemaController.FindByCode)
	router.POST("/api/cinema", cinemaController.Create)
	router.PUT("/api/cinema/:companyCode", cinemaController.Update)
	router.DELETE("/api/categories/:companyCode", cinemaController.Delete)

	router.PanicHandler = exception.ErrorHandler

	return router
}

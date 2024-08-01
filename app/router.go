package app

import (
	"net/http"
	"x1-cinema/controller"
	"x1-cinema/exception"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(cinemaController controller.CinemaController) *httprouter.Router {
	router := httprouter.New()
	apiPrefix := "/api/v1"

	router.GET(apiPrefix+"/get", cinemaController.FindAll)
	router.GET(apiPrefix+"/get/:CinemaCode", cinemaController.FindByCode)
	router.POST(apiPrefix+"/create", cinemaController.Create)
	router.PUT(apiPrefix+"/update/:CinemaCode", cinemaController.Update)
	router.DELETE(apiPrefix+"/delete/:CinemaCode", cinemaController.Delete)

	router.ServeFiles("/swagger/*filepath", http.Dir("./swagger"))

	router.PanicHandler = exception.ErrorHandler

	return router
}

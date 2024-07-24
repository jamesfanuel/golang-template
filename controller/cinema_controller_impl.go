package controller

import (
	"encoding/json"
	"net/http"
	"x1-cinema/helper"
	"x1-cinema/model/web"
	"x1-cinema/service"

	"github.com/julienschmidt/httprouter"
)

type CinemaControllerImpl struct {
	CinemaService service.CinemaService
}

func NewCinemaController(cinemaService service.CinemaService) CinemaController {
	return &CinemaControllerImpl{
		CinemaService: cinemaService,
	}
}

func (controller *CinemaControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	decoder := json.NewDecoder(request.Body)

	cinemaCreateRequest := web.CinemaCreateRequest{}
	err := decoder.Decode(&cinemaCreateRequest)
	helper.PanicIfError(err)

	cinemaResponse := controller.CinemaService.Create(request.Context(), cinemaCreateRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   cinemaResponse,
	}

	encoder := json.NewEncoder(writer)
	err = encoder.Encode(webResponse)
	helper.PanicIfError(err)
}

func (controller *CinemaControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	decoder := json.NewDecoder(request.Body)

	cinemaUpdateRequest := web.CinemaUpdateRequest{}
	err := decoder.Decode(&cinemaUpdateRequest)
	helper.PanicIfError(err)

	cinemaResponse := controller.CinemaService.Update(request.Context(), cinemaUpdateRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   cinemaResponse,
	}

	encoder := json.NewEncoder(writer)
	err = encoder.Encode(webResponse)
	helper.PanicIfError(err)
}

func (controller *CinemaControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	controller.CinemaService.Delete(request.Context(), params.ByName("cinema_code"))
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(webResponse)
	helper.PanicIfError(err)
}

func (controller *CinemaControllerImpl) FindByCode(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	cinemaResponse := controller.CinemaService.FindByCode(request.Context(), params.ByName("cinema_code"))
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   cinemaResponse,
	}

	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(webResponse)
	helper.PanicIfError(err)
}

func (controller *CinemaControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	cinemaResponses := controller.CinemaService.FindAll(request.Context())
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   cinemaResponses,
	}

	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(webResponse)
	helper.PanicIfError(err)
}

package controller

import (
	"encoding/json"
	"fmt"
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

// @Summary Create a new theater
// @Description Create a new theater with the provided data.
// @Tags cinema
// @Accept json
// @Produce json
// @Param Cinema body web.CinemaCreateRequest true "Cinema object to be created"
// @Success 200 {object} web.WebResponse
// @Router /create [post]
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

// @Summary Update theater
// @Description Update theater with the provided data.
// @Tags cinema
// @Accept json
// @Produce json
// @Param CinemaCode path string true "Cinema object to be updated"
// @Param Cinema body web.CinemaUpdateRequest true "Cinema object to be updated"
// @Success 200 {object} web.WebResponse
// @Router /update/{CinemaCode} [put]
func (controller *CinemaControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	decoder := json.NewDecoder(request.Body)

	cinemaUpdateRequest := web.CinemaUpdateRequest{}
	err := decoder.Decode(&cinemaUpdateRequest)
	helper.PanicIfError(err)

	fmt.Print(params.ByName("CinemaCode"))

	cinemaResponse := controller.CinemaService.Update(request.Context(), cinemaUpdateRequest, params.ByName("CinemaCode"))
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   cinemaResponse,
	}

	encoder := json.NewEncoder(writer)
	err = encoder.Encode(webResponse)
	helper.PanicIfError(err)
}

// @Summary Delete theater
// @Description Delete theater with Cinema Code provided.
// @Tags cinema
// @Accept json
// @Produce json
// @Param CinemaCode path string true "Cinema object to be deleted"
// @Success 200 {object} web.WebResponse
// @Router /delete/{CinemaCode} [delete]
func (controller *CinemaControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	controller.CinemaService.Delete(request.Context(), params.ByName("CinemaCode"))
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(webResponse)
	helper.PanicIfError(err)
}

// @Summary Find By Theater Code
// @Description Find Specific Theater By Code Provided.
// @Tags cinema
// @Accept json
// @Produce json
// @Param CinemaCode path string true "Cinema Found"
// @Success 200 {object} web.WebResponse
// @Router /get/{CinemaCode} [get]
func (controller *CinemaControllerImpl) FindByCode(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	cinemaResponse := controller.CinemaService.FindByCode(request.Context(), params.ByName("CinemaCode"))
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

// @Summary Find All Theater
// @Description Find All Theater on Master.
// @Tags cinema
// @Accept json
// @Produce json
// @Success 200 {object} web.WebResponse
// @Router /get [get]
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

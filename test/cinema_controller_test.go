package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"x1-cinema/app"
	"x1-cinema/controller"
	"x1-cinema/helper"
	"x1-cinema/middleware"
	"x1-cinema/model/domain"
	"x1-cinema/repository"
	"x1-cinema/service"

	"github.com/go-playground/validator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/magiconair/properties/assert"
)

func setupTestDB() *sql.DB {
	db, err := sql.Open("mysql", "dsserver:xxi2121.@tcp(k8s.devel.intra.db.cinema21.co.id:3306)/db_xone")
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func setupRouter(db *sql.DB) http.Handler {
	validate := validator.New()
	cinemaRepository := repository.NewCinemaRepository()
	cinemaService := service.NewCinemaService(cinemaRepository, db, validate)
	cinemaController := controller.NewCinemaController(cinemaService)
	router := app.NewRouter(cinemaController)

	return middleware.NewAuthMiddleware(router)
}

func truncateCinema(db *sql.DB) {
	db.Exec("TRUNCATE mg_cinema")
}

func TestCreateCinemaSuccess(t *testing.T) {
	db := setupTestDB()
	truncateCinema(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"cinema_code" : "Cinema", "cinema_name" : "Cinema"}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:6010/api/v1/create", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	// assert.Equal(t, "Gadget", responseBody["data"].(map[string]interface{})["cinema_name"])
}

func TestCreateCinemaFailed(t *testing.T) {
	db := setupTestDB()
	truncateCinema(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"cinema_code" : "", "cinema_name" : ""}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:6010/api/v1/create", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 500, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 500, int(responseBody["code"].(float64)))
	assert.Equal(t, "INTERNAL SERVER ERROR", responseBody["status"])
}

func TestUpdateCinemaSuccess(t *testing.T) {
	db := setupTestDB()
	truncateCinema(db)

	tx, _ := db.Begin()
	cinemaRepository := repository.NewCinemaRepository()
	cinema := cinemaRepository.Save(context.Background(), tx, domain.Cinema{
		CinemaCode: "Cinema",
		CinemaName: "Cinema",
	})
	tx.Commit()

	router := setupRouter(db)

	requestBody := strings.NewReader(`{"cinema_name" : "Cinema"}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:6010/api/v1/update/"+cinema.CinemaCode, requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, cinema.CinemaCode, responseBody["data"].(map[string]interface{})["cinema_code"])
	assert.Equal(t, "Cinema", responseBody["data"].(map[string]interface{})["cinema_name"])
}

func TestUpdateCinemaFailed(t *testing.T) {
	db := setupTestDB()
	truncateCinema(db)

	tx, _ := db.Begin()
	cinemaRepository := repository.NewCinemaRepository()
	cinema := cinemaRepository.Save(context.Background(), tx, domain.Cinema{
		CinemaCode: "Cinema",
		CinemaName: "Cinema",
	})
	tx.Commit()

	router := setupRouter(db)

	requestBody := strings.NewReader(`{"cinema_name" : "-"}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:6010/api/cinema/"+cinema.CinemaCode, requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	// assert.Equal(t, 404, int(responseBody["code"].(float64)))
	// assert.Equal(t, "INTERNAL SERVER ERROR", responseBody["status"])
}

func TestGetCinemaSuccess(t *testing.T) {
	db := setupTestDB()
	truncateCinema(db)

	tx, _ := db.Begin()
	cinemaRepository := repository.NewCinemaRepository()
	cinema := cinemaRepository.Save(context.Background(), tx, domain.Cinema{
		CinemaCode: "Cinema",
		CinemaName: "Cinema",
	})
	tx.Commit()

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:6010/api/v1/get/"+cinema.CinemaCode, nil)
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, cinema.CinemaCode, responseBody["data"].(map[string]interface{})["cinema_code"])
	assert.Equal(t, cinema.CinemaName, responseBody["data"].(map[string]interface{})["cinema_name"])
}

func TestGetCinemaFailed(t *testing.T) {
	db := setupTestDB()
	truncateCinema(db)
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:6010/api/v1/get/nothing", nil)
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 500, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	// assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}

func TestDeleteCinemaSuccess(t *testing.T) {
	db := setupTestDB()
	truncateCinema(db)

	tx, _ := db.Begin()
	cinemaRepository := repository.NewCinemaRepository()
	cinema := cinemaRepository.Save(context.Background(), tx, domain.Cinema{
		CinemaCode: "Cinema",
		CinemaName: "Cinema",
	})
	tx.Commit()

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:6010/api/cinema/"+cinema.CinemaCode, nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	// assert.Equal(t, 200, int(responseBody["code"].(float64)))
	// assert.Equal(t, "OK", responseBody["status"])
}

func TestDeleteCinemaFailed(t *testing.T) {
	db := setupTestDB()
	truncateCinema(db)
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:6010/api/cinema/404", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}

func TestListcinemaSuccess(t *testing.T) {
	db := setupTestDB()
	truncateCinema(db)

	tx, _ := db.Begin()
	cinemaRepository := repository.NewCinemaRepository()
	cinema1 := cinemaRepository.Save(context.Background(), tx, domain.Cinema{
		CinemaCode: "Cinema1",
		CinemaName: "Cinema1",
	})
	cinema2 := cinemaRepository.Save(context.Background(), tx, domain.Cinema{
		CinemaCode: "Cinema2",
		CinemaName: "Cinema2",
	})
	tx.Commit()

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:6010/api/v1/get", nil)
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])

	fmt.Println(responseBody)

	var cinema = responseBody["data"].([]interface{})

	cinemaResponse1 := cinema[0].(map[string]interface{})
	cinemaResponse2 := cinema[1].(map[string]interface{})

	assert.Equal(t, cinema1.CinemaCode, cinemaResponse1["cinema_code"])
	assert.Equal(t, cinema1.CinemaName, cinemaResponse1["cinema_name"])

	assert.Equal(t, cinema2.CinemaCode, cinemaResponse2["cinema_code"])
	assert.Equal(t, cinema2.CinemaName, cinemaResponse2["cinema_name"])
}

func TestUnauthorized(t *testing.T) {
	db := setupTestDB()
	// truncateCinema(db)
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:6010/api/cinema", nil)
	request.Header.Add("X-API-Key", "SALAH")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 401, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 401, int(responseBody["code"].(float64)))
	assert.Equal(t, "UNAUTHORIZED", responseBody["status"])
}

package test

import (
	"context"
	"encoding/json"
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
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func OpenConnection() *gorm.DB {
	dialect := mysql.Open("dsserver:xxi2121.@tcp(k8s.devel.intra.db.cinema21.co.id:3306)/db_xone?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(dialect, &gorm.Config{})

	helper.PanicIfError(err)

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	return db
}

var db = OpenConnection()

func TestOpenConnection(t *testing.T) {
	assert.NotNil(t, db)
}

func setupRouter(db *gorm.DB) http.Handler {
	validate := validator.New()
	cinemaRepository := repository.NewCinemaRepository(db)
	cinemaService := service.NewCinemaService(cinemaRepository, db, validate)
	cinemaController := controller.NewCinemaController(cinemaService)
	router := app.NewRouter(cinemaController)

	return middleware.NewAuthMiddleware(router)
}

func deleteCinema(db *gorm.DB, params string) {
	db.Exec("DELETE FROM mg_cinema WHERE cinema_code = " + params)
}

func TestCreateCinemaSuccess(t *testing.T) {
	db := OpenConnection()
	deleteCinema(db, "CINETEST")
	router := setupRouter(db)

	requestBody := strings.NewReader(`
		{
			"cinema_code":      "CINETEST",
			"cinema_name":      "CINETEST123",
			"cinema_owner":     "OWNER123",
			"location_code":    "LOC123",
			"province_code":    "PROV123",
			"city_code":        "CITY123",
			"region_code":      "REGION123",
			"company_code":     "COMOWN123",
			"cinema_level":     "REG",
			"oracle_code":      "ORACLECODE",
			"is_data_migration": "N",
			"close_flag":       "N",
			"close_start":      "2024-01-01",
			"close_end":        "2024-01-01",
			"operator_email":   "domain@co.id",
			"created_by":       "USER"
		}
	`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:6010/api/v1/create", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	// assert.Equal(t, requestBody, responseBody["data"])
	// assert.Equal(t, "Gadget", responseBody["data"].(map[string]interface{})["name"])
}

// func TestCreateCategoryFailed(t *testing.T) {
// 	db := setupTestDB()
// 	truncateCategory(db)
// 	router := setupRouter(db)

// 	requestBody := strings.NewReader(`{"name" : ""}`)
// 	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/categories", requestBody)
// 	request.Header.Add("Content-Type", "application/json")
// 	request.Header.Add("X-API-Key", "RAHASIA")

// 	recorder := httptest.NewRecorder()

// 	router.ServeHTTP(recorder, request)

// 	response := recorder.Result()
// 	assert.Equal(t, 400, response.StatusCode)

// 	body, _ := io.ReadAll(response.Body)
// 	var responseBody map[string]interface{}
// 	json.Unmarshal(body, &responseBody)

// 	assert.Equal(t, 400, int(responseBody["code"].(float64)))
// 	assert.Equal(t, "BAD REQUEST", responseBody["status"])
// }

func TestUpdateCinemaSuccess(t *testing.T) {
	db := OpenConnection()
	deleteCinema(db, "TESTUPDATE")

	tx := db.Begin()
	cinemaRepository := repository.NewCinemaRepository(db)
	cinema, _ := cinemaRepository.Save(context.Background(), tx, domain.Cinema{
		CinemaCode: "TESTUPDATE",
		CinemaName: "TESTUPDATE",
	})
	tx.Commit()

	router := setupRouter(db)

	requestBody := strings.NewReader(`
		{
			"cinema_name":      "CINETEST123",
			"cinema_owner":     "OWNER123",
			"location_code":    "LOC123",
			"province_code":    "PROV123",
			"city_code":        "CITY123",
			"region_code":      "REGION123",
			"company_code":     "COMOWN123",
			"cinema_level":     "REG",
			"oracle_code":      "ORACLECODE",
			"is_data_migration": "N",
			"close_flag":       "N",
			"close_start":      "2024-01-01",
			"close_end":        "2024-01-01",
			"operator_email":   "domain@co.id",
			"updated_by":       "USER"
		}
	`)

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
}

// func TestUpdateCategoryFailed(t *testing.T) {
// 	db := setupTestDB()
// 	truncateCategory(db)

// 	tx, _ := db.Begin()
// 	categoryRepository := repository.NewCategoryRepository()
// 	category := categoryRepository.Save(context.Background(), tx, domain.Category{
// 		Name: "Gadget",
// 	})
// 	tx.Commit()

// 	router := setupRouter(db)

// 	requestBody := strings.NewReader(`{"name" : ""}`)
// 	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), requestBody)
// 	request.Header.Add("Content-Type", "application/json")
// 	request.Header.Add("X-API-Key", "RAHASIA")

// 	recorder := httptest.NewRecorder()

// 	router.ServeHTTP(recorder, request)

// 	response := recorder.Result()
// 	assert.Equal(t, 400, response.StatusCode)

// 	body, _ := io.ReadAll(response.Body)
// 	var responseBody map[string]interface{}
// 	json.Unmarshal(body, &responseBody)

// 	assert.Equal(t, 400, int(responseBody["code"].(float64)))
// 	assert.Equal(t, "BAD REQUEST", responseBody["status"])
// }

func TestGetCinemaSuccess(t *testing.T) {
	db := OpenConnection()
	deleteCinema(db, "TESTGET")

	tx := db.Begin()
	cinemaRepository := repository.NewCinemaRepository(db)
	cinema, _ := cinemaRepository.Save(context.Background(), tx, domain.Cinema{
		CinemaCode: "TESTGET",
		CinemaName: "TESTGET",
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
	// assert.Equal(t, category.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, cinema.CinemaCode, responseBody["data"].(map[string]interface{})["cinema_code"])
}

// func TestGetCategoryFailed(t *testing.T) {
// 	db := setupTestDB()
// 	truncateCategory(db)
// 	router := setupRouter(db)

// 	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories/404", nil)
// 	request.Header.Add("X-API-Key", "RAHASIA")

// 	recorder := httptest.NewRecorder()

// 	router.ServeHTTP(recorder, request)

// 	response := recorder.Result()
// 	assert.Equal(t, 404, response.StatusCode)

// 	body, _ := io.ReadAll(response.Body)
// 	var responseBody map[string]interface{}
// 	json.Unmarshal(body, &responseBody)

// 	assert.Equal(t, 404, int(responseBody["code"].(float64)))
// 	assert.Equal(t, "NOT FOUND", responseBody["status"])
// }

func TestDeleteCinemaSuccess(t *testing.T) {
	db := OpenConnection()
	deleteCinema(db, "TESTDELETE")

	tx := db.Begin()
	cinemaRepository := repository.NewCinemaRepository(db)
	cinema, _ := cinemaRepository.Save(context.Background(), tx, domain.Cinema{
		CinemaCode: "TESTDELETE",
		CinemaName: "TESTDELETE",
	})
	tx.Commit()

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/v1/delete/"+cinema.CinemaCode, nil)
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
}

// func TestDeleteCategoryFailed(t *testing.T) {
// 	db := setupTestDB()
// 	truncateCategory(db)
// 	router := setupRouter(db)

// 	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/categories/404", nil)
// 	request.Header.Add("Content-Type", "application/json")
// 	request.Header.Add("X-API-Key", "RAHASIA")

// 	recorder := httptest.NewRecorder()

// 	router.ServeHTTP(recorder, request)

// 	response := recorder.Result()
// 	assert.Equal(t, 404, response.StatusCode)

// 	body, _ := io.ReadAll(response.Body)
// 	var responseBody map[string]interface{}
// 	json.Unmarshal(body, &responseBody)

// 	assert.Equal(t, 404, int(responseBody["code"].(float64)))
// 	assert.Equal(t, "NOT FOUND", responseBody["status"])
// }

func TestListCinemasSuccess(t *testing.T) {
	db := OpenConnection()
	deleteCinema(db, "TESTLIST1")
	deleteCinema(db, "TESTLIST2")

	tx := db.Begin()
	cinemaRepository := repository.NewCinemaRepository(db)
	cinema1, _ := cinemaRepository.Save(context.Background(), tx, domain.Cinema{
		CinemaCode: "TESTLIST1",
		CinemaName: "TESTLIST1",
	})
	cinema2, _ := cinemaRepository.Save(context.Background(), tx, domain.Cinema{
		CinemaCode: "TESTLIST2",
		CinemaName: "TESTLIST2",
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

	var cinema = responseBody["data"].([]interface{})

	cinemaResponse1 := cinema[0].(map[string]interface{})
	cinemaResponse2 := cinema[1].(map[string]interface{})

	assert.Equal(t, cinema1.CinemaCode, cinemaResponse1["cinema_code"])
	assert.Equal(t, cinema1.CinemaName, cinemaResponse1["cinema_name"])

	assert.Equal(t, cinema2.CinemaCode, cinemaResponse2["cinema_code"])
	assert.Equal(t, cinema2.CinemaName, cinemaResponse2["cinema_name"])
}

func TestUnauthorized(t *testing.T) {
	db := OpenConnection()
	// truncateCategory(db)
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:6010/api/v1/get", nil)
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

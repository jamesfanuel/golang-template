package test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"x1-cinema/app"
	"x1-cinema/controller"
	"x1-cinema/helper"
	"x1-cinema/middleware"
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

func TestCreateCinemaSuccess(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"cinema_code" : "test", "cinema_name" : "test"}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/cinema", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)
}

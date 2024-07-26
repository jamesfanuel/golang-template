package test

import (
	"testing"
	"time"
	"x1-cinema/helper"
	"x1-cinema/model/domain"

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

func TestCreateCinema(t *testing.T) {
	// id := uuid.New()
	user := domain.Cinema{
		CinemaCode:   "TESTCINEMACODE",
		CinemaName:   "TESTCINEMANAME",
		ProvinceCode: "TESTPROVINCECODE",
		CityCode:     "TESTCITYCODE",
		RegionCode:   "TESTREGIONCODE",
		CinemaLevel:  "TESTCINEMALEVEL",
		// RowId:        id,
	}

	response := db.Create(&user)
	assert.Nil(t, response.Error)
	assert.Equal(t, int64(1), response.RowsAffected)
}

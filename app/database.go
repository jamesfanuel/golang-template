package app

import (
	"database/sql"
	"time"
	"x1-cinema/helper"
)

func NewDB() *sql.DB {
	db, err := sql.Open("mysql", "dsserver:xxi2121.@tcp(k8s.devel.intra.db.cinema21.co.id:3306)/db_xone")
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

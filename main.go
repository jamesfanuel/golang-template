package main

import (
	"fmt"
	"net/http"
	"x1-cinema/helper"
	"x1-cinema/middleware"

	// "x1-cinema/helper"

	_ "github.com/go-sql-driver/mysql"
)

func NewServer(authMiddleware *middleware.AuthMiddleware) *http.Server {
	return &http.Server{
		Addr:    ":6010",
		Handler: authMiddleware,
		// Handler: router,
	}
}

// @title x1-cinema API
// @version 1.0
// @description API Doc for x1-cinema.
// @BasePath /api/v1
func main() {
	// app.NewLog("Info", "Initiate Application")
	// app.NewEureka()

	fmt.Print("Ready To Serve\n")

	server := InitializedServer()
	err := server.ListenAndServe()

	helper.PanicIfError(err)
}

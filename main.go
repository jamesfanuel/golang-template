package main

import (
	"fmt"
	"go-ms-template-service/helper"
	"go-ms-template-service/middleware"
	"net/http"

	// "go-ms-template-service/helper"

	_ "github.com/go-sql-driver/mysql"
)

func NewServer(authMiddleware *middleware.AuthMiddleware) *http.Server {
	return &http.Server{
		Addr:    ":6010",
		Handler: authMiddleware,
		// Handler: router,
	}
}

// @title go-ms-template-service API
// @version 1.0
// @description API Doc for go-ms-template-service.
// @BasePath /api/v1
func main() {
	// app.NewLog("Info", "Initiate Application")
	// app.NewEureka()

	fmt.Print("Ready To Serve\n")

	server := InitializedServer()
	err := server.ListenAndServe()

	helper.PanicIfError(err)
}

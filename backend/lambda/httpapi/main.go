package main

import (
	"net/http"
	
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Recover())
	
	e.GET("/", handleHello)
	
	//nolint:errcheck
	e.Start(":12001")
}

func handleHello(c echo.Context) error {
	return c.String(http.StatusOK, "hello, "+c.RealIP())
}

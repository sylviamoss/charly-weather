package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	router := echo.New()
	router.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Charly Weather is up")
	})

	weatherModule := NewModule()
	weatherModule.RegisterRoutes(router)

	router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	router.Logger.Fatal(router.Start(":" + os.Getenv("PORT")))
}

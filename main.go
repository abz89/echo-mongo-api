package main

import (
	"echo-mongo-api/configs"
	"echo-mongo-api/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// run database
	configs.ConnectDB()

	// middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// routes
	routes.UserRoute(e)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, &echo.Map{"message": "Hello World!"})
	})

	e.Logger.Fatal(e.Start(":" + configs.GoDotEnvVarible("PORT")))
}

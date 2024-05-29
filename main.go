package main

import (
	"echo-mongo-api/configs"
	"echo-mongo-api/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Connect database
	configs.ConnectDB()

	// Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	routes.AuthRoute(e)
	routes.UserRoute(e)

	// Test route
	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, &echo.Map{"message": "Hello World!"})
	})

	// Start web server
	e.Logger.Fatal(e.Start(":" + configs.GoDotEnvVariable("PORT")))
}

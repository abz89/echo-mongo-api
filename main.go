package main

import (
	"echo-mongo-api/configs"
	"echo-mongo-api/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// run database
	configs.ConnectDB()

	// routes
	routes.UserRoute(e)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, &echo.Map{"message": "Hello World!"})
	})

	e.Logger.Fatal(e.Start(":" + configs.GoDotEnvVarible("PORT")))
}

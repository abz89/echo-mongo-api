package main

import (
	"echo-mongo-api/configs"
	"echo-mongo-api/routes"

	"github.com/labstack/echo/v4"
)

// Adapted from this article https://dev.to/hackmamba/build-a-rest-api-with-golang-and-mongodb-echo-version-2gdg (https://github.com/Mr-Malomz/echo-mongo-api)
func main() {
	e := echo.New()

	// run database
	configs.ConnectDB()

	// routes
	routes.UserRoute(e)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, &echo.Map{"message": "Hello World!"})
	})

	e.Logger.Fatal(e.Start(":8080"))
}

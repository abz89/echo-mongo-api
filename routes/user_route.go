package routes

import (
	"echo-mongo-api/controllers"

	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Echo) {
	e.POST("/users", controllers.CreateUser)
	e.GET("/users", controllers.ListUsers)
	e.GET("/users/:userId", controllers.GetUser)
	e.PUT("/users/:userId", controllers.EditUser)
	e.PATCH("/users/:userId", controllers.PatchUser)
	e.DELETE("/users/:userId", controllers.DeleteUser)
}

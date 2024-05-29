package routes

import (
	"echo-mongo-api/configs"
	"echo-mongo-api/controllers"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// User routes
func UserRoute(e *echo.Echo) {
	e.POST("/users", controllers.CreateUser)
	e.GET("/users", controllers.ListUsers, middleware.BasicAuth(checkBasicAuth))
	e.GET("/users/:userId", controllers.GetUser, JWTmiddleware(), checkUserIdFromJWT)
	e.PUT("/users/:userId", controllers.EditUser, JWTmiddleware(), checkUserIdFromJWT)
	e.PATCH("/users/:userId", controllers.PatchUser, JWTmiddleware(), checkUserIdFromJWT)
	e.DELETE("/users/:userId", controllers.DeleteUser, JWTmiddleware(), checkUserIdFromJWT)
}

// function for basic auth
func checkBasicAuth(username, password string, _ echo.Context) (bool, error) {
	if username == configs.GoDotEnvVariable("ADMIN_USER") && password == configs.GoDotEnvVariable("ADMIN_PASSWORD") {
		return true, nil
	}

	return false, nil
}

// function for JWT middleware
func JWTmiddleware() echo.MiddlewareFunc {
	config := echojwt.Config{NewClaimsFunc: func(c echo.Context) jwt.Claims {
		return new(jwtCustomClaims)
	},
		SigningKey: []byte(configs.GoDotEnvVariable("SECRET")),
	}

	return echojwt.WithConfig(config)
}

// function for checking current user
func checkUserIdFromJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*jwtCustomClaims)
		userId := c.Param("userId")
		if userId != claims.Id {
			return echo.ErrUnauthorized
		}
		return next(c)
	}

}

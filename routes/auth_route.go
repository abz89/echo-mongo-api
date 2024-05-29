package routes

import (
	"context"
	"echo-mongo-api/configs"
	"echo-mongo-api/controllers"
	"echo-mongo-api/models"
	"time"

	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"golang.org/x/crypto/bcrypt"
)

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type jwtCustomClaims struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Admin    bool   `json:"admin"`
	jwt.RegisteredClaims
}

// Auth routes
func AuthRoute(e *echo.Echo) {
	e.POST("/login", login)
	e.POST("/register", controllers.CreateUser)
}

// Login controller
func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := models.UserCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)

	if err != nil {
		return echo.ErrUnauthorized
	}

	if !verifyPassword(user.Id, password) {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &jwtCustomClaims{
		user.Id.Hex(),
		user.Username,
		user.Name,
		user.Email,
		user.Admin,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(configs.GoDotEnvVariable("SECRET")))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"token": t, "user": user})
}

// Verify password for current user
func verifyPassword(userId primitive.ObjectID, password string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var credentials struct {
		Id       primitive.ObjectID `bson:"_id"`
		Password string             `bson:"password"`
		Username string             `bson:"username"`
	}

	objecId, _ := primitive.ObjectIDFromHex(userId.Hex())

	err := models.UserCollection.FindOne(ctx, bson.M{"_id": objecId}).Decode(&credentials)

	if err != nil {
		return false
	}

	return checkPasswordHash(password, credentials.Password)
}

// Check if the password hash is correct
func checkPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

package controllers

import (
	"context"
	"echo-mongo-api/models"
	"echo-mongo-api/responses"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Instance of the validator library
var validate = validator.New()

// Create a new user
func CreateUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	// validate the request body
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
	}

	// use the validator library to validate required fields
	if validateErr := validate.Struct(&user); validateErr != nil {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: validateErr.Error()})
	}

	hashedPassword, _ := models.HashPassword(user.Password)

	newUser := models.User{
		Id:       primitive.NewObjectID(),
		Username: user.Username,
		Email:    user.Email,
		Password: hashedPassword,
		Name:     user.Name,
		Admin:    user.Admin,
		Title:    user.Title,
		Location: user.Location,
	}

	result, err := models.UserCollection.InsertOne(ctx, newUser)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	var createdUser models.User

	err2 := models.UserCollection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&createdUser)

	if err2 != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err2.Error()})
	}

	return c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: createdUser})
}

// List all users
func ListUsers(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get query parameters for pagination
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	skip, err := strconv.Atoi(c.QueryParam("skip"))
	if err != nil || skip < 0 {
		skip = 0
	}

	// Find users with limit and skip
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(skip))

	results, err := models.UserCollection.Find(ctx, bson.M{}, findOptions)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	// reading from the db in an optimal way
	defer results.Close(ctx)
	var users []models.User
	for results.Next(ctx) {
		var user models.User
		if err = results.Decode(&user); err != nil {
			return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		}

		users = append(users, user)
	}

	// Get total count of users
	total, err := models.UserCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	response := responses.PaginatedResponse{
		Total: int(total),
		Limit: limit,
		Skip:  skip,
		Data:  users,
	}

	return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: response})
}

// Get a single user
func GetUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Param("userId")
	var user models.User
	defer cancel()

	objectId, _ := primitive.ObjectIDFromHex(userId)

	err := models.UserCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&user)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: user})
}

// Edit user using PUT method
func EditUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Param("userId")
	var user models.User
	defer cancel()

	objectId, _ := primitive.ObjectIDFromHex(userId)

	// validate the request body
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
	}

	// use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
	}

	hashedPassword, _ := models.HashPassword(user.Password)

	update := bson.M{"username": user.Username, "email": user.Email, "password": hashedPassword, "name": user.Name, "admin": user.Admin, "title": user.Title, "location": user.Location}

	result, err := models.UserCollection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$set": update})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	// get updated user details
	var updatedUser models.User
	if result.MatchedCount == 1 {
		err := models.UserCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&updatedUser)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		}
	} else {
		return c.JSON(http.StatusNotFound, responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: fmt.Sprintf("user with id %s not found", userId)})
	}

	return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: updatedUser})
}

// Edit user using PATCH method
func PatchUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Param("userId")
	var user models.User
	defer cancel()

	objectId, _ := primitive.ObjectIDFromHex(userId)

	err := models.UserCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&user)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	var patch map[string]interface{}
	if err := c.Bind(&patch); err != nil {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
	}

	// Merge the fields from the request into the existing document
	for key, value := range patch {
		switch key {
		case "username":
			user.Username = value.(string)
		case "email":
			user.Email = value.(string)
		case "password":
			hashedPassword, _ := models.HashPassword(value.(string))
			user.Password = hashedPassword
		case "name":
			user.Name = value.(string)
		case "admin":
			user.Admin = value.(bool)
		case "title":
			user.Title = value.(string)
		case "location":
			user.Location = value.(string)
		}
	}

	result, err := models.UserCollection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$set": user})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	// get patched user details
	var patchedUser models.User
	if result.MatchedCount == 1 {
		err := models.UserCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&patchedUser)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		}
	}

	return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: patchedUser})
}

// Delete a user
func DeleteUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Param("userId")
	defer cancel()

	ObjectId, _ := primitive.ObjectIDFromHex(userId)

	result, err := models.UserCollection.DeleteOne(ctx, bson.M{"_id": ObjectId})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	if result.DeletedCount < 1 {
		return c.JSON(http.StatusNotFound, responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: "user not found"})
	}

	return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: "user deleted"})
}

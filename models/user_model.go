package models

import (
	"echo-mongo-api/configs"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// User struct
type User struct {
	Id       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username,omitempty" validate:"required"`
	Email    string             `json:"email,omitempty" validate:"required,email"`
	Password string             `json:"password,omitempty" validate:"required"`
	Name     string             `json:"name,omitempty" validate:"required"`
	Admin    bool               `json:"admin,omitempty"`
	Title    string             `json:"title,omitempty"`
	Location string             `json:"location,omitempty"`
}

// User collection
var UserCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

// Hash password before save using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

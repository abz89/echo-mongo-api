package configs

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect to MongoDB
func ConnectDB() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(GoDotEnvVariable("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}

	return client
}

// Client instance
var DB = ConnectDB()

// Getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	dbName, err := extractDBName(GoDotEnvVariable("MONGO_URI"))
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database(dbName).Collection(collectionName)
	return collection
}

// Function to extract the database name from the URI
func extractDBName(uri string) (string, error) {
	if !strings.HasPrefix(uri, "mongodb://") && !strings.HasPrefix(uri, "mongodb+srv://") {
		return "", fmt.Errorf("invalid MongoDB URI: %s", uri)
	}

	// Remove the protocol part
	parts := strings.SplitN(uri, "/", 4)
	if len(parts) < 4 {
		return "", fmt.Errorf("database name not found in URI: %s", uri)
	}

	dbName := parts[3]
	if dbName == "" {
		return "", fmt.Errorf("database name is empty in URI: %s", uri)
	}

	return dbName, nil
}

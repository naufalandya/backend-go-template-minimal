package mongodb

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func InitMongoDB() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI is not set in the environment... (；´д｀)")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Test connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Cannot ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB~ nya~! ✨🐱")
	MongoClient = client
}

func GetCollection(databaseName, collectionName string) *mongo.Collection {
	return MongoClient.Database(databaseName).Collection(collectionName)
}

package firebase

import (
	"context"
	"fmt"
	"log"
	"modular_monolith/server/functions"
	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var App *firebase.App

func InitFirebase() {
	err := LoadEnv() // Just in case .env isn't loaded yet
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Get the credentials file path from environment
	credsFile := functions.LoadEnvVariable("FIREBASE_CREDENTIALS", "firebase_credentials.json")

	opt := option.WithCredentialsFile(credsFile)
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase App: %v", err)
	}

	App = app
	fmt.Println("Connected to Firebase~ 🐣✨")
}

func LoadEnv() error {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		return nil // ignore if no env
	}
	return nil
}

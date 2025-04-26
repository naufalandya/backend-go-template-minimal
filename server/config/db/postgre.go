package db

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

var DB *pgx.Conn

func loadEnvVariable(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func InitDB() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := loadEnvVariable("DB_USER", "postgres")
	dbPassword := loadEnvVariable("DB_PASSWORD", "password-default")
	dbHost := loadEnvVariable("DB_HOST", "localhost")
	dbPort := loadEnvVariable("DB_PORT", "5432")
	dbName := loadEnvVariable("DB_NAME", "postgres")
	dbSchema := loadEnvVariable("DB_SCHEMA", "public")
	env := loadEnvVariable("ENV", "development")
	sslMode := loadEnvVariable("SSL_STATUS_DB", "disable")

	fmt.Println("DB_USER:", dbUser)
	fmt.Println("DB_PASSWORD:", dbPassword)
	fmt.Println("DB_HOST:", dbHost)
	fmt.Println("DB_PORT:", dbPort)
	fmt.Println("DB_NAME:", dbName)
	fmt.Println("DB_SCHEMA:", dbSchema)
	fmt.Println("ENV:", env)

	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(dbUser, dbPassword),
		Host:   fmt.Sprintf("%s:%s", dbHost, dbPort),
		Path:   dbName,
	}
	q := dsn.Query()
	q.Add("sslmode", sslMode)
	if dbSchema != "" {
		q.Add("search_path", dbSchema)
	}
	dsn.RawQuery = q.Encode()

	databaseURL := dsn.String()

	DB, err = pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	fmt.Println("Connected to PostgreSQL~ ✨")
}

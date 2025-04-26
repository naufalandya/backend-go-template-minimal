package db

import (
	"context"
	"fmt"
	"log"
	"modular_monolith/server/functions"
	"net/url"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

var DB *pgx.Conn

func InitDB() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := functions.LoadEnvVariable("DB_USER", "postgres")
	dbPassword := functions.LoadEnvVariable("DB_PASSWORD", "password-default")
	dbHost := functions.LoadEnvVariable("DB_HOST", "localhost")
	dbPort := functions.LoadEnvVariable("DB_PORT", "5432")
	dbName := functions.LoadEnvVariable("DB_NAME", "postgres")
	dbSchema := functions.LoadEnvVariable("DB_SCHEMA", "public")
	env := functions.LoadEnvVariable("ENV", "development")
	sslMode := functions.LoadEnvVariable("SSL_STATUS_DB", "disable")

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

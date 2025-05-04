package db

import (
	"context"
	"fmt"
	"log"
	"modular_monolith/server/functions"
	"net/url"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB     *pgx.Conn // for raw SQL kawaii ✨
	GormDB *gorm.DB  // for ORM kawaii ✨
)

func InitDB() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file, nyaaa (´；ω；`)!")
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

	// Create DSN (kawaii string ✨)
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

	// Connect using pgx (raw SQL mode 🧋)
	DB, err = pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatal("Unable to connect to database with pgx, sobs: ", err)
	}
	fmt.Println("Connected to PostgreSQL (pgx)~!! (๑˃̵ᴗ˂̵)و")

	// Connect using GORM (ORM mode 🧋)
	// GORM logger level based on ENV (production = silent, dev = info)
	var gormLogger logger.Interface
	if env == "production" {
		gormLogger = logger.Default.LogMode(logger.Silent)
	} else {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	GormDB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  databaseURL,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatal("Unable to connect to database with GORM, sobs: ", err)
	}

	// Test GORM connection with a lil ping-ping~ 💬
	sqlDB, err := GormDB.DB()
	if err != nil {
		log.Fatal("Failed to get generic DB from GORM: ", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database from GORM: ", err)
	}
	fmt.Println("Connected to PostgreSQL (GORM)~!! (＾▽＾)")
}

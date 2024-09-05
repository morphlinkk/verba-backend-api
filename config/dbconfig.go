package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// DB содержит соединение с БД
var DB *pgxpool.Pool

// InitDB Создает подключение к БД по данным из .env
func InitDB() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	dbConfig := DatabaseConfig{
		User:     os.Getenv("DB_USER"),
		Host:     os.Getenv("DB_HOST"),
		Database: os.Getenv("DB_DATABASE"),
		Password: os.Getenv("DB_PASSWORD"),
		Port:     os.Getenv("DB_PORT"),
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		fmt.Printf("Unable to parse config: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	DB, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	err = DB.Ping(ctx)
	if err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")
}

// DatabaseConfig Представляет конфигурацию базы данных
type DatabaseConfig struct {
	User     string
	Host     string
	Database string
	Password string
	Port     string
}

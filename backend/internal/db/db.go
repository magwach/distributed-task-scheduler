package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
)

func Connect() (*pgxpool.Pool, error) {
	databaseUrl := os.Getenv("DATABASE_URL")

	if databaseUrl == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	config, err := pgxpool.ParseConfig(databaseUrl)

	if err != nil {
		log.Fatalf("Unable to parse DATABASE_URL: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	log.Println("Database connection pool established")

	return pool, nil
}

func Close() {
	if pool != nil {
		pool.Close()
		log.Println("Database connection pool closed")
	}
}

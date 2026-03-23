package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
)

func Connect(databaseUrl string) (*pgxpool.Pool, error) {

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

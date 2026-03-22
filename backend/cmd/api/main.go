package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/magwach/distributed-task-scheduler/backend/internal/db"
)

func main() {

	err := godotenv.Load("../../.env")

	if err != nil {
		log.Println("Warning: No env file found")
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	dbUrl := os.Getenv("DATABASE_URL")

	if dbUrl == "" {
		log.Fatal("DATABASE_URL is not set")
	}


	pool := db.Connect()

	if pool != nil {
		log.Fatalf("Unable to connect to DB, %v", err)
	}

	defer db.Close()

	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	log.Printf("Server is running on port: %v", port)

	if err = app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

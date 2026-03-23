package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/magwach/distributed-task-scheduler/backend/internal/db"
	"github.com/magwach/distributed-task-scheduler/backend/internal/routes"
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

	pool, err := db.Connect()

	if err != nil {
		log.Fatalf("Unable to connect to DB, %v", err)
	}

	defer db.Close()

	app := fiber.New()

	c := cors.New(
		cors.Config{
			AllowOrigins: "http://localhost:3030",
			AllowHeaders: "Content-Type, Accept, Authorization",
			AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		},
	)

	app.Use(c)

	v1Routes := app.Group("/api/v1")

	v1Routes.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	taskRoutes := routes.NewTaskRoutes(v1Routes, pool)

	taskRoutes.TaskRoutes()

	if err = app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	log.Printf("Server is running on port: %v", port)

}

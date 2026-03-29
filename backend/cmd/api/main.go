package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/magwach/distributed-task-scheduler/backend/internal/auth"
	"github.com/magwach/distributed-task-scheduler/backend/internal/db"
	"github.com/magwach/distributed-task-scheduler/backend/internal/models"
	"github.com/magwach/distributed-task-scheduler/backend/internal/queue"
	"github.com/magwach/distributed-task-scheduler/backend/internal/routes"
	"github.com/magwach/distributed-task-scheduler/backend/internal/websockets"
)

func main() {

	err := godotenv.Load("../../.env")

	if err != nil {
		log.Println("Warning: No env file found")
	}

	client := os.Getenv("CLIENT_URL")

	if client == "" {
		log.Fatal("CLIENT_URL is not set")
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	redisUrl := os.Getenv("REDIS_ADDR")

	if redisUrl == "" {
		log.Fatal("REDIS_ADDR is not set")
	}

	dbUrl := os.Getenv("DATABASE_URL")

	if dbUrl == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	pool, err := db.Connect(dbUrl)

	if err != nil {
		log.Fatalf("Unable to connect to DB, %v", err)
	}

	defer db.Close()

	app := fiber.New()

	c := cors.New(
		cors.Config{
			AllowOrigins:     client,
			AllowCredentials: true,
			AllowHeaders:     "Content-Type, Accept, Authorization",
			AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		},
	)

	app.Use(c)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	v1Routes := app.Group("/api/v1", auth.AuthMiddleware)

	queue.InitRedis(redisUrl)

	hub := websockets.HubInit()

	go func() {
		sub := queue.GetRedisClient().Subscribe(context.Background(), "task:updates")
		ch := sub.Channel()
		for msg := range ch {
			var event models.TaskUpdateEvent
			err := json.Unmarshal([]byte(msg.Payload), &event)
			if err != nil {
				log.Println("Failed to unmarshal task update:", err)
				continue
			}
			hub.Broadcast(event)
		}
	}()

	websocketRoutes := routes.NewWebSocketRoutes(v1Routes, hub)
	taskRoutes := routes.NewTaskRoutes(v1Routes, pool)
	authRoutes := routes.NewAuthRoutes(app, pool)
	userRoutes := routes.NewUserRoutes(v1Routes)

	taskRoutes.TaskRoutes()
	websocketRoutes.WebSocketRoutes()
	authRoutes.AuthRoutes()
	userRoutes.UserRoutes()

	if err = app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	log.Printf("Server is running on port: %v", port)

}

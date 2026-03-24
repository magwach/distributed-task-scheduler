package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/magwach/distributed-task-scheduler/backend/internal/db"
	"github.com/magwach/distributed-task-scheduler/backend/internal/services"
)

func main() {

	err := godotenv.Load("../../.env")

	if err != nil {
		log.Println("Warning: No env file found")
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

	scheduler := services.SchedulerServiceImpl(pool, redisUrl)

	go func() {
		ticker := time.NewTicker(5 * time.Second)

		defer ticker.Stop()

		for range ticker.C {
			scheduler.ProcessPendingTasks()
		}

	}()

	select {}

}

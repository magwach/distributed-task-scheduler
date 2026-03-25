package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/magwach/distributed-task-scheduler/backend/internal/db"
	"github.com/magwach/distributed-task-scheduler/backend/internal/models"
	"github.com/magwach/distributed-task-scheduler/backend/internal/queue"
	"github.com/magwach/distributed-task-scheduler/backend/internal/retry"
	"github.com/magwach/distributed-task-scheduler/backend/internal/services"
	"github.com/redis/go-redis/v9"
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

	DB, err := db.Connect(dbUrl)

	if err != nil {
		log.Fatalf("Unable to connect to DB, %v", err)
	}

	defer DB.Close()

	queue.InitRedis(redisUrl)

	updateTaskExcecutionStatusToSuccessQuery := `
		UPDATE task_excecutions
		SET status = 'success', finished_at = now()
		WHERE task_id = $1
		`

	updateTaskStatusToSuccessQuery := `
		UPDATE tasks
		SET status = 'success'
		WHERE id = $1
		`

	updateTaskExcecutionStatusToFailedQuery := `
		UPDATE task_excecutions
		SET status = 'failed', finished_at = now(), error_message = $1
		WHERE task_id = $2
		`

	updateTaskStatusToFailedQuery := `
		UPDATE tasks
		SET status = 'failed'
		WHERE id = $1
		`
	getTaskDetailsQuery := `
	SELECT *
	FROM tasks
	WHERE id = $1
	`

	incrementTheRetriesQuery := `
	UPDATE tasks
	SET status = 'pending', next_run_at = now() + $1, retry_count = $2
	WHERE id = $3
	`
	getExcecutionIdQuery := `
	SELECT id
	FROM task_excecutions
	WHERE task_id = $1
	`

	for {
		task := models.Task{}
		taskId, err := queue.Dequeue(2 * time.Second)

		if err != nil {
			if err == redis.Nil {
				continue
			}
			log.Println("Redis error:", err)
			continue
		}

		err = DB.QueryRow(context.Background(),
			getTaskDetailsQuery,
			taskId,
		).Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Schedule,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.NextRunAt,
			&task.LastRunAt,
			&task.MaxRetries,
			&task.RetryCount,
			&task.RetryDelaySeconds,
		)

		if err != nil {
			log.Println("Failed to fetch task details for ID:", taskId)
			continue
		}

		go func(task models.Task) {

			var taskExcecutionId string

			err := DB.QueryRow(context.Background(),
				getExcecutionIdQuery,
				task.ID,
			).Scan(
				&taskExcecutionId,
			)

			if err != nil {
				log.Println("Failed to get task excecutin Id for task: ", task.ID)
				return
			}

			err = services.WriteLog(context.Background(), DB, taskExcecutionId, "info", fmt.Sprintf("Worker picked up task %v (id: %v)", task.Title, task.ID))

			if err != nil {
				return
			}

			startTime, err := services.WorkerFunction(task, context.Background(), DB, taskExcecutionId)

			if err != nil {
				if task.RetryCount < task.MaxRetries {

					delay := retry.Delay(task.RetryDelaySeconds, task.RetryCount)

					_, err = DB.Exec(context.Background(),
						incrementTheRetriesQuery,
						delay,
						task.RetryCount+1,
						task.ID)

					if err != nil {
						log.Println("Failed to increment retry count: ", err)
						return
					}

					err = services.WriteLog(context.Background(), DB, taskExcecutionId, "warn", fmt.Sprintf("Task failed. Scheduling retry %v/%v in %vs", task.RetryCount, task.MaxRetries, delay))

					if err != nil {
						return
					}
				} else {
					err = services.WriteLog(context.Background(), DB, taskExcecutionId, "error", fmt.Sprintf("Task permanently failed after %v attempts", task.MaxRetries))

					if err != nil {
						return
					}
				}

				_, err = DB.Exec(context.Background(),
					updateTaskExcecutionStatusToFailedQuery,
					err,
					task.ID,
				)
				if err != nil {
					log.Println("Failed to update task execution:", err)
					return
				}

				_, err = DB.Exec(context.Background(),
					updateTaskStatusToFailedQuery,
					task.ID,
				)
				if err != nil {
					log.Println("Failed to update task :", err)
					return
				}

				err = services.WriteLog(context.Background(), DB, taskExcecutionId, "error", fmt.Sprintf("Task failed: %v", err))

				if err != nil {
					return
				}

			} else {
				_, err = DB.Exec(context.Background(),
					updateTaskExcecutionStatusToSuccessQuery,
					task.ID,
				)
				if err != nil {
					log.Println("Failed to update task execution:", err)
					return
				}
				_, err = DB.Exec(context.Background(),
					updateTaskStatusToSuccessQuery,
					task.ID,
				)
				if err != nil {
					log.Println("Failed to update task :", err)
					return
				}
				duration := time.Since(startTime)
				err = services.WriteLog(context.Background(), DB, taskExcecutionId, "info", fmt.Sprintf("Task completed successfully in %v ms", duration.Milliseconds()))

				if err != nil {
					return
				}
			}
		}(task)
	}

}

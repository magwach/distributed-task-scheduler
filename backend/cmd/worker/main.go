package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/magwach/distributed-task-scheduler/backend/internal/db"
	"github.com/magwach/distributed-task-scheduler/backend/internal/models"
	"github.com/magwach/distributed-task-scheduler/backend/internal/queue"
	"github.com/magwach/distributed-task-scheduler/backend/internal/retry"
	"github.com/magwach/distributed-task-scheduler/backend/internal/services"
	"github.com/magwach/distributed-task-scheduler/backend/pkg/utils"
	"github.com/redis/go-redis/v9"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found, using environment variables")
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
		SET status = 'success', next_run_at = $1, last_run_at = now(), updated_at = now()
		WHERE id = $2
		`

	updateTaskExcecutionStatusToFailedQuery := `
		UPDATE task_excecutions
		SET status = 'failed', finished_at = now(), error_message = $1
		WHERE task_id = $2
		`

	updateTaskStatusToFailedQuery := `
		UPDATE tasks
		SET status = 'failed', retry_count = 0
		WHERE id = $1
		`
	getTaskDetailsQuery := `
	SELECT *
	FROM tasks
	WHERE id = $1
	`

	incrementTheRetriesQuery := `
	UPDATE tasks
	SET status = 'pending', next_run_at = $1, retry_count = $2, retry_delay_seconds = $3
	WHERE id = $4
	`
	getExcecutionIdQuery := `
	SELECT id
	FROM task_excecutions
	WHERE task_id = $1
	ORDER BY created_at DESC
	LIMIT 1
	`

	go func() {
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
				&task.Priority,
			)

			if err != nil {
				log.Println("Failed to fetch task details for ID:", taskId, " error: ", err)
				continue
			}

			go func(task models.Task) {

				var taskExcecutionId string
				var nextRunAt time.Time

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

				startTime, workerErr := services.WorkerFunction(task, context.Background(), DB, taskExcecutionId)

				if workerErr != nil {
					if task.RetryCount < task.MaxRetries {

						retryDelaySeconds := retry.Delay(task.RetryDelaySeconds, task.RetryCount)

						delay := time.Now().Add(time.Duration(retryDelaySeconds) * time.Second)
						_, err = DB.Exec(context.Background(),
							incrementTheRetriesQuery,
							delay,
							task.RetryCount+1,
							retryDelaySeconds,
							task.ID)

						if err != nil {
							log.Println("Failed to increment retry count: ", err)
							return
						}

						err = services.WriteLog(context.Background(), DB, taskExcecutionId, "warning", fmt.Sprintf("Task failed. Scheduling retry %v/%v in %vs", task.RetryCount+1, task.MaxRetries, retryDelaySeconds))

						if err != nil {
							return
						}

						return
					} else {
						err = services.WriteLog(context.Background(), DB, taskExcecutionId, "error", fmt.Sprintf("Task permanently failed after %v attempts", task.MaxRetries))

						if err != nil {
							return
						}
					}

					_, err = DB.Exec(context.Background(),
						updateTaskExcecutionStatusToFailedQuery,
						workerErr.Error(),
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

					errMsg := workerErr.Error()

					updateEvent := models.TaskUpdateEvent{
						TaskID:       task.ID,
						Status:       "failed",
						UpdatedAt:    time.Now(),
						ExecutionID:  taskExcecutionId,
						ErrorMessage: &errMsg,
						RetryCount:   &task.RetryCount,
						MaxRetries:   &task.MaxRetries,
						NextRunAt:    &nextRunAt,
					}

					data, err := json.Marshal(updateEvent)

					if err != nil {
						log.Println("Failed to parse the message to JSON")
						return
					}

					err = queue.GetRedisClient().Publish(context.Background(), "task:updates", data).Err()
					if err != nil {
						log.Println("Failed to publish task update:", err)
					}

					err = services.WriteLog(context.Background(), DB, taskExcecutionId, "error", fmt.Sprintf("Task failed: %v", errMsg))

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

					nextRunAt, err = utils.ParseCron(task.Schedule)

					if err != nil {
						log.Println("Failed to parse cron :", err)
						return
					}

					_, err = DB.Exec(context.Background(),
						updateTaskStatusToSuccessQuery,
						nextRunAt,
						task.ID,
					)

					if err != nil {
						log.Println("Failed to update task :", err)
						return
					}

					updateEvent := models.TaskUpdateEvent{
						TaskID:      task.ID,
						Status:      "success",
						UpdatedAt:   time.Now(),
						ExecutionID: taskExcecutionId,
						NextRunAt:   &nextRunAt,
					}

					data, err := json.Marshal(updateEvent)

					if err != nil {
						log.Println("Failed to parse the message to JSON")
						return
					}

					err = queue.GetRedisClient().Publish(context.Background(), "task:updates", data).Err()
					if err != nil {
						log.Println("Failed to publish task update:", err)
					}

					duration := time.Since(startTime)
					err = services.WriteLog(context.Background(), DB, taskExcecutionId, "info", fmt.Sprintf("Task completed successfully in %v ms", duration.Milliseconds()))

					if err != nil {
						return
					}
				}
			}(task)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Worker service running"))
	})

	log.Printf("Starting dummy web server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}

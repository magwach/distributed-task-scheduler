package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magwach/distributed-task-scheduler/backend/internal/lock"
	"github.com/magwach/distributed-task-scheduler/backend/internal/models"
	"github.com/magwach/distributed-task-scheduler/backend/internal/queue"
)

type schedulerService struct {
	DB *pgxpool.Pool
}

func SchedulerServiceImpl(db *pgxpool.Pool) *schedulerService {
	return &schedulerService{
		DB: db,
	}
}

func (s *schedulerService) ProcessPendingTasks() {

	tasks := []models.Task{}

	getAllTasksWithPendingStatusQuery := `
	SELECT *
	FROM tasks
	WHERE next_run_at <= now()
	AND status != 'running'
	`

	rows, err := s.DB.Query(context.Background(),
		getAllTasksWithPendingStatusQuery,
	)

	if err != nil {
		log.Println("Failed to get pending tasks: ", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		task := models.Task{}
		err := rows.Scan(
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
			log.Println("Failed to scan task:", err)
			continue
		}

		tasks = append(tasks, task)
	}

	if len(tasks) == 0 {
		return
	}

	updateTasksToRunningQuery := `
	UPDATE tasks
	SET status = 'running'
	WHERE id = $1 AND status = 'pending'
	`

	createTasksExcecutionRecordQuery := `
	INSERT INTO task_excecutions (task_id, status, started_at)
	VALUES ($1, 'running', now())
	RETURNING id
	`

	for _, task := range tasks {

		lockAquired, lockValue := lock.AquireLock(task.ID)

		if !lockAquired {
			log.Println("Task already locked, skipping:", task.ID)
			continue
		}

		go func(t models.Task, lockValue string) {
			var executionID string

			defer lock.ReleaseLock(t.ID, lockValue)

			_, err := s.DB.Exec(context.Background(),
				updateTasksToRunningQuery,
				t.ID,
			)

			if err != nil {
				log.Println("Failed to update the status")
				return
			}

			err = s.DB.QueryRow(context.Background(),
				createTasksExcecutionRecordQuery,
				t.ID,
			).Scan(&executionID)

			if err != nil {
				log.Println("Failed to insert execution record:", err)
				return
			}

			err = queue.Enqueue(t.ID)

			if err != nil {
				log.Println("Failed to add the task: ", t.Title, " to redis.")
				return
			}

			updateEvent := models.TaskUpdateEvent{
				TaskID:      t.ID,
				Status:      "running",
				UpdatedAt:   time.Now(),
				ExecutionID: executionID,
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

		}(task, lockValue)
	}

	log.Println("Processed", len(tasks), "tasks")

}

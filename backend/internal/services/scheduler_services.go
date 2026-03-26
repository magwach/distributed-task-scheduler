package services

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magwach/distributed-task-scheduler/backend/internal/models"
	"github.com/magwach/distributed-task-scheduler/backend/internal/queue"
	"github.com/magwach/distributed-task-scheduler/backend/internal/websockets"
)

type schedulerService struct {
	DB *pgxpool.Pool
	WB *websockets.Hub
}

func SchedulerServiceImpl(db *pgxpool.Pool, wb *websockets.Hub) *schedulerService {
	return &schedulerService{
		DB: db,
		WB: wb,
	}
}

func (s *schedulerService) ProcessPendingTasks() {

	tasks := []models.Task{}

	var executionID string

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
	WHERE id = $1
	`

	createTasksExcecutionRecordQuery := `
	INSERT INTO task_excecutions (task_id, status, started_at)
	VALUES ($1, 'running', now())
	RETURNING id
	`

	for _, task := range tasks {

		go func(t models.Task) {
			_, err := s.DB.Exec(context.Background(),
				updateTasksToRunningQuery,
				task.ID,
			)

			if err != nil {
				log.Println("Failed to update the status")
				return
			}

			err = s.DB.QueryRow(context.Background(),
				createTasksExcecutionRecordQuery,
				task.ID,
			).Scan(&executionID)

			if err != nil {
				log.Println("Failed to insert execution record:", err)
				return
			}

			err = queue.Enqueue(task.ID)

			if err != nil {
				log.Println("Failed to add the task: ", task.Title, " to redis.")
				return
			}

			updateEvent := models.TaskUpdateEvent{
				TaskID:      task.ID,
				Status:      "running",
				UpdatedAt:   time.Now(),
				ExecutionID: executionID,
			}

			s.WB.Broadcast(updateEvent)
		}(task)
	}

	log.Println("Processed", len(tasks), "tasks")

}

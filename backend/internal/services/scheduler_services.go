package services

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
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
	WHERE (next_run_at <= now() OR next_run_at IS NULL)
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
			return
		}
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
	RETURNING *
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

			_, err = s.DB.Exec(context.Background(),
				createTasksExcecutionRecordQuery,
				task.ID,
			)

			if err != nil {
				log.Println("Failed to insert execution record:", err)
				return
			}

			err = queue.Enqueue(task.ID)

			if err != nil {
				log.Println("Failed to add the task: ", task.Title, " to redis.")
				return
			}

		}(task)

	}

	log.Println("Processed", len(tasks), "tasks")

}

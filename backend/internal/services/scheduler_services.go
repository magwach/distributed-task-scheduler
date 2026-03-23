package services

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magwach/distributed-task-scheduler/backend/internal/models"
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
	WHERE status = 'pending'
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
		)

		if err != nil {
			log.Println("Failed to scan task:", err)
			return
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

	for _, task := range tasks {
		_, err := s.DB.Exec(context.Background(),
			updateTasksToRunningQuery,
			task.ID,
		)

		if err != nil {
			log.Fatal("Failed to update the status")
			return
		}
	}

	createTasksExcecutionRecordQuery := `
	INSERT INTO task_excecutions (task_id, status, started_at)
	VALUES ($1, 'running', now())
	RETURNING *
	`

	for _, value := range tasks {
		_, err := s.DB.Exec(context.Background(),
			createTasksExcecutionRecordQuery,
			value.ID,
		)

		if err != nil {
			log.Println("Failed to insert execution record:", err)
			return
		}
	}

	log.Println("Processed", len(tasks), "tasks")

	for _, task := range tasks {

		task := task

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

		go func(task models.Task) {
			err := WorkerFunction(task)

			if err != nil {
				_, err = s.DB.Exec(context.Background(),
					updateTaskExcecutionStatusToFailedQuery,
					err,
					task.ID,
				)
				if err != nil {
					log.Println("Failed to update task execution:", err)
					return
				}

				_, err = s.DB.Exec(context.Background(),
					updateTaskStatusToFailedQuery,
					task.ID,
				)
				if err != nil {
					log.Println("Failed to update task :", err)
					return
				}
			} else {
				_, err = s.DB.Exec(context.Background(),
					updateTaskExcecutionStatusToSuccessQuery,
					task.ID,
				)
				if err != nil {
					log.Println("Failed to update task execution:", err)
					return
				}
				_, err = s.DB.Exec(context.Background(),
					updateTaskStatusToSuccessQuery,
					task.ID,
				)
				if err != nil {
					log.Println("Failed to update task :", err)
					return
				}

			}
		}(task)
	}

}

package services

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magwach/distributed-task-scheduler/backend/internal/models"
	"github.com/magwach/distributed-task-scheduler/backend/internal/queue"
)

type schedulerService struct {
	DB       *pgxpool.Pool
	RedisUrl string
}

func SchedulerServiceImpl(db *pgxpool.Pool, redisUrl string) *schedulerService {
	return &schedulerService{
		DB:       db,
		RedisUrl: redisUrl,
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

	queue.InitRedis(s.RedisUrl)

	for _, task := range tasks {
		err := queue.Enqueue(task.ID)

		if err != nil {
			log.Println("Failed to add the task: ", task.Title, " to redis.")
		}
	}

	log.Println("Processed", len(tasks), "tasks")



}

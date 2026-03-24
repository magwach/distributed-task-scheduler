package services

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magwach/distributed-task-scheduler/backend/internal/models"
	"github.com/magwach/distributed-task-scheduler/backend/internal/queue"
	"github.com/robfig/cron/v3"
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

	updateTaskNextRunTimeQuery := `
	UPDATE tasks 
	SET next_run_at = $1
	WHERE id = $2
	`

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

		parser := cron.NewParser(
			cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
		)

		schedule, err := parser.Parse(task.Schedule)

		if err != nil {
			log.Println("Failed to parse the cron")
			return
		}

		nextRun := schedule.Next(time.Now())

		task.NextRunAt = nextRun

		_, err = s.DB.Exec(context.Background(),
			updateTaskNextRunTimeQuery,
			nextRun,
			task.ID,
		)
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

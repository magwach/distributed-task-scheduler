package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magwach/distributed-task-scheduler/backend/internal/dto"
	"github.com/magwach/distributed-task-scheduler/backend/internal/models"
	"github.com/robfig/cron/v3"
)

type TaskService struct {
	DB *pgxpool.Pool
}

func NewTaskService(db *pgxpool.Pool) TaskService {
	return TaskService{
		DB: db,
	}
}

func (s *TaskService) CreateTask(taskInput dto.CreateTaskRequest) (*models.Task, error) {

	task := models.Task{}

	query := `
	INSERT INTO tasks (title, description, schedule, next_run_at)
	VALUES ($1, $2, $3, $4)
	RETURNING id, title, description, schedule, status, created_at, updated_at
	`

	parser := cron.NewParser(
		cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
	)

	schedule, err := parser.Parse(taskInput.Schedule)

	if err != nil {
		log.Println("Failed to parse the cron")
		return nil, fmt.Errorf("invalid cron expression: %w", err)
	}

	nextRun := schedule.Next(time.Now())

	err = s.DB.QueryRow(
		context.Background(),
		query,
		taskInput.Title,
		taskInput.Description,
		taskInput.Schedule,
		nextRun,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Schedule,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *TaskService) GetTasks() ([]models.Task, error) {
	tasks := []models.Task{}

	query := `
	SELECT id, title, description, schedule, status, created_at, updated_at
	FROM tasks
	`
	rows, err := s.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TaskService) GetTask(id uuid.UUID) (*models.Task, error) {
	task := models.Task{}

	query := `
	SELECT id, title, description, schedule, status, created_at, updated_at
	FROM tasks
	WHERE id = $1
	`

	err := s.DB.QueryRow(context.Background(),
		query,
		id,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Schedule,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("Task not found")
		}
		return nil, err
	}

	return &task, nil
}

func (s *TaskService) DeleteTask(id uuid.UUID) (*models.Task, error) {
	task := models.Task{}

	query := `
	DELETE FROM tasks
	WHERE id = $1
	RETURNING id, title, description, schedule, status, created_at, updated_at
	`
	err := s.DB.QueryRow(context.Background(),
		query,
		id,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Schedule,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("Task not found")
		}
		return nil, err
	}

	return &task, nil

}

func (s *TaskService) GetLogs(taskId uuid.UUID) ([]models.TaskLog, error) {

	logs := []models.TaskLog{}
	var executionId string

	getLogsQuery := `
	SELECT *
	FROM task_logs
	WHERE execution_id = $1
	`

	getTaskExcecutionId := `
	SELECT id
	FROM task_excecutions
	WHERE task_id = $1
	`

	err := s.DB.QueryRow(context.Background(),
		getTaskExcecutionId,
		taskId,
	).Scan(
		&executionId,
	)

	rows, err := s.DB.Query(context.Background(),
		getLogsQuery,
		executionId,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		log := models.TaskLog{}

		rows.Scan(
			&log.ID,
			&log.ExecutionID,
			&log.Level,
			&log.Message,
			&log.CreatedAt,
		)
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

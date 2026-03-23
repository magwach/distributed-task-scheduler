package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magwach/distributed-task-scheduler/backend/internal/dto"
	"github.com/magwach/distributed-task-scheduler/backend/internal/models"
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
	INSERT INTO tasks (title, description, schedule)
	VALUES ($1, $2, $3)
	RETURNING *
	`

	err := s.DB.QueryRow(
		context.Background(),
		query,
		taskInput.Title,
		taskInput.Description,
		taskInput.Schedule,
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
	SELECT *
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
	SELECT *
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
	RETURNING *
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

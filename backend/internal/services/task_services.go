package services

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magwach/distributed-task-scheduler/backend/internal/dto"
	"github.com/magwach/distributed-task-scheduler/backend/internal/models"
	"github.com/magwach/distributed-task-scheduler/backend/pkg/utils"
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

	nextRun, err := utils.ParseCron(taskInput.Schedule)

	if err != nil {
		return nil, err
	}

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
	executions := []models.TaskExecution{}

	getTaskQuery := `
	SELECT id, title, description, schedule, status, created_at, updated_at, next_run_at, last_run_at, max_retries, retry_count, retry_delay_seconds 
	FROM tasks
	WHERE id = $1
	`

	getExcecutionsQuery := `
	SELECT *
	FROM task_excecutions
	WHERE task_id = $1
	`

	getLogsQuery := `
	SELECT *
	FROM task_logs
	WHERE execution_id = $1
	`

	err := s.DB.QueryRow(context.Background(),
		getTaskQuery,
		id,
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("Task not found")
		}
		return nil, err
	}

	excecutionRows, err := s.DB.Query(context.Background(),
		getExcecutionsQuery,
		id,
	)

	if err != nil {
		return nil, err
	}

	defer excecutionRows.Close()

	for excecutionRows.Next() {
		execution := models.TaskExecution{}

		err := excecutionRows.Scan(
			&execution.ID,
			&execution.TaskID,
			&execution.Status,
			&execution.StartedAt,
			&execution.FinishedAt,
			&execution.ErrorMessage,
			&execution.CreatedAt,
			&execution.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		executions = append(executions, execution)

	}

	if err := excecutionRows.Err(); err != nil {
		return nil, err
	}

	task.Executions = executions

	for index, execution := range executions {

		logs := []models.TaskLog{}

		logRows, err := s.DB.Query(context.Background(),
			getLogsQuery,
			execution.ID,
		)

		if err != nil {
			return nil, err
		}

		for logRows.Next() {
			log := models.TaskLog{}

			err := logRows.Scan(
				&log.ID,
				&log.ExecutionID,
				&log.Level,
				&log.Message,
				&log.CreatedAt,
			)

			if err != nil {
				logRows.Close()
				return nil, err
			}

			logs = append(logs, log)
		}

		if err := logRows.Err(); err != nil {
			logRows.Close()
			return nil, err
		}

		logRows.Close()

		for index, taskExecution := range task.Executions {
			if taskExecution.ID == execution.ID {
				task.Executions[index].Logs = logs
			}
			continue
		}

		task.Executions[index].Logs = logs

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
	ORDER BY created_at DESC
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

func (s *TaskService) RetryTask(taskId uuid.UUID) error {
	query := `
	UPDATE tasks
	SET status = 'pending', retry_count = 0, retry_delay_seconds = 60, next_run_at = now()
	WHERE id = $1
	RETURNING id, title, description, schedule, status, created_at, updated_at 
	`

	_, err := s.DB.Exec(context.Background(),
		query,
		taskId,
	)

	if err != nil {
		log.Println("Failed to retry: ", err)
		return err
	}

	return nil

}

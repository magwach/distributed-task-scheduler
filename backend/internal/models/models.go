package models

import (
	"time"
)

type Task struct {
	ID                string    `json:"id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	Schedule          string    `json:"schedule"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	NextRunAt         time.Time `json:"next_run_at"`
	LastRunAt         time.Time `json:"last_run_at"`
	MaxRetries        int       `json:"max_retries"`
	RetryCount        int       `json:"retry_count"`
	RetryDelaySeconds int       `json:"retry_delay_seconds"`
}

type TaskExecution struct {
	ID           string     `json:"id"`
	TaskID       string     `json:"task_id"`
	Status       string     `json:"status"`
	StartedAt    time.Time  `json:"started_at"`
	FinishedAt   *time.Time `json:"finished_at"`
	ErrorMessage string     `json:"error_message"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type TaskLog struct {
	ID          string    `json:"id"`
	ExecutionID string    `json:"task_id"`
	Level       string    `json:"level"`
	Message     string    `json:"message"`
	CreatedAt   time.Time `json:"created_at"`
}

type TaskUpdateEvent struct {
	TaskID       string     `json:"task_id"`
	Status       string     `json:"status"`
	UpdatedAt    time.Time  `json:"updated_at"`
	ExecutionID  string     `json:"execution_id"`
	NextRunAt    *time.Time `json:"next_run_at,omitempty"`
	ErrorMessage *string    `json:"error_message,omitempty"`
	RetryCount   *int       `json:"retry_count,omitempty"`
	MaxRetries   *int       `json:"max_retries,omitempty"`
}

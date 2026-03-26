package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magwach/distributed-task-scheduler/backend/internal/models"
)

func WorkerFunction(task models.Task, ctx context.Context, db *pgxpool.Pool, executionID string) (time.Time, error) {
	log.Println("Excecuting task: ", task.Title)
	startTime := time.Now()

	err := WriteLog(ctx, db, executionID, "info", "Execution started")

	if err != nil {
		return time.Time{}, errors.New("test error")

	}

	time.Sleep(3 * time.Second)

	log.Println("Finished excecuting task: ", task.Title)
	return startTime, nil
}

func WriteLog(ctx context.Context, db *pgxpool.Pool, executionID string, level string, message string) error {

	writeToLogTableQuery := `
	INSERT INTO task_logs (execution_id, level, message, created_at)
	VALUES ($1, $2, $3, now())
	RETURNING *
	`
	_, err := db.Exec(ctx,
		writeToLogTableQuery,
		executionID,
		level,
		message,
	)

	if err != nil {
		log.Println("Error trying to insert log in the table: ", err)
		return errors.New("failed to insert log to table")
	}
	return nil
}

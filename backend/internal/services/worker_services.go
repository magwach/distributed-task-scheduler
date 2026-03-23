package services

import (
	"log"
	"time"

	"github.com/magwach/distributed-task-scheduler/backend/internal/models"
)

func WorkerFunction(task models.Task) error {

	log.Println("Excecuting task: ", task.Title)

	time.Sleep(3 * time.Second)

	log.Println("Finished excecuting task: ", task.Title)
	return nil
}

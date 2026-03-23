package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magwach/distributed-task-scheduler/backend/internal/handlers"
	"github.com/magwach/distributed-task-scheduler/backend/internal/services"
)

type TaskRoutesImpl struct {
	App fiber.Router
	DB  *pgxpool.Pool
}

func NewTaskRoutes(app fiber.Router, db *pgxpool.Pool) *TaskRoutesImpl {
	return &TaskRoutesImpl{
		App: app,
		DB:  db,
	}
}

func (r *TaskRoutesImpl) TaskRoutes() {
	taskService := services.NewTaskService(r.DB)

	taskHandlers := handlers.NewTaskHandler(&taskService)

	r.App.Post("/tasks", taskHandlers.CreateTask)
	r.App.Get("/tasks", taskHandlers.GetTasks)
	r.App.Get("/tasks/:id", taskHandlers.GetTask)
	r.App.Delete("/tasks/:id", taskHandlers.DeleteTask)
}

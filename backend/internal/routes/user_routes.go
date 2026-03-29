package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/magwach/distributed-task-scheduler/backend/internal/handlers"
)

type UserRoutesImpl struct {
	App fiber.Router
}

func NewUserRoutes(app fiber.Router) *UserRoutesImpl {
	return &UserRoutesImpl{
		App: app,
	}
}

func (r *UserRoutesImpl) UserRoutes() {
	userHandlers := handlers.NewUserHandler()

	r.App.Get("/me", userHandlers.Me)
}

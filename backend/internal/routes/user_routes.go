package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magwach/distributed-task-scheduler/backend/internal/handlers"
	"github.com/magwach/distributed-task-scheduler/backend/internal/services"
)

type UserRoutesImpl struct {
	App fiber.Router
	DB  *pgxpool.Pool
}

func NewUserRoutes(app fiber.Router, db *pgxpool.Pool) *UserRoutesImpl {
	return &UserRoutesImpl{
		App: app,
		DB:  db,
	}
}

func (r *UserRoutesImpl) UserRoutes() {
	userService := services.NewUserService(r.DB)
	userHandlers := handlers.NewUserHandler(&userService)

	r.App.Get("/me", userHandlers.Me)
}

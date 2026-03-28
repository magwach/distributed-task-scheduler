package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magwach/distributed-task-scheduler/backend/internal/handlers"
	"github.com/magwach/distributed-task-scheduler/backend/internal/services"
)

type AuthRoutesImpl struct {
	App fiber.Router
	DB  *pgxpool.Pool
}

func NewAuthRoutes(app fiber.Router, db *pgxpool.Pool) *AuthRoutesImpl {
	return &AuthRoutesImpl{
		App: app,
		DB:  db,
	}
}

func (r *TaskRoutesImpl) AuthRoutes() {
	authService := services.NewAuthService(r.DB)

	authHandlers := handlers.NewAuthHandler(&authService)

	r.App.Post("/auth/register", authHandlers.Register)
	r.App.Post("/auth/login", authHandlers.Login)
	r.App.Post("/auth/refresh", authHandlers.Refresh)
	r.App.Get("/auth/google", handleIt)
	r.App.Get("/auth/google/callback", handleIt)
	r.App.Get("/auth/github", handleIt)
	r.App.Get("/auth/github/callback", handleIt)

}

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

func (r *AuthRoutesImpl) AuthRoutes() {
	authService := services.NewAuthService(r.DB)

	authHandlers := handlers.NewAuthHandler(&authService)

	r.App.Post("/auth/register", authHandlers.Register)
	r.App.Post("/auth/login", authHandlers.Login)
	r.App.Post("/auth/refresh", authHandlers.Refresh)
	r.App.Post("/auth/logout", authHandlers.Logout)
	r.App.Get("/auth/google", authHandlers.GoogleLogin)
	r.App.Get("/auth/google/callback", authHandlers.GoogleCallback)
	r.App.Get("/auth/github", authHandlers.GitHubLogin)
	r.App.Get("/auth/github/callback", authHandlers.GitHubCallback)
}

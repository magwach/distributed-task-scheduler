package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/magwach/distributed-task-scheduler/backend/internal/handlers"
	"github.com/magwach/distributed-task-scheduler/backend/internal/websockets"
)

type WebSocketRoutesImpl struct {
	App fiber.Router
}

func NewWebSocketRoutes(app fiber.Router) *WebSocketRoutesImpl {
	return &WebSocketRoutesImpl{
		App: app,
	}
}

func (r *WebSocketRoutesImpl) WebSocketRoutes() {
	hub := websockets.HubInit()

	webSocketHandler := handlers.NewWebSocketHandler(hub)

	r.App.Get("/ws", webSocketHandler.RegisterRoutes)

}

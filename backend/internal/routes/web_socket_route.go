package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/magwach/distributed-task-scheduler/backend/internal/handlers"
	"github.com/magwach/distributed-task-scheduler/backend/internal/websockets"
)

type WebSocketRoutesImpl struct {
	App fiber.Router
	Hub *websockets.Hub
}

func NewWebSocketRoutes(app fiber.Router, hub *websockets.Hub) *WebSocketRoutesImpl {
	return &WebSocketRoutesImpl{
		App: app,
		Hub: hub,
	}
}

func (r *WebSocketRoutesImpl) WebSocketRoutes() {

	webSocketHandler := handlers.NewWebSocketHandler(r.Hub)

	r.App.Get("/ws", webSocketHandler.RegisterRoutes)

}

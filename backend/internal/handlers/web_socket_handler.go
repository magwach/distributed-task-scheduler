package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/magwach/distributed-task-scheduler/backend/internal/websockets"
)

type WebSocketHanlderImpl struct {
	Hub *websockets.Hub
}

func NewWebSocketHandler(hub *websockets.Hub) *WebSocketHanlderImpl {
	return &WebSocketHanlderImpl{
		Hub: hub,
	}
}

func (h *WebSocketHanlderImpl) RegisterRoutes(c *fiber.Ctx) error {
	return websocket.New(func(conn *websocket.Conn) {
		h.Hub.AddClient(conn)
		defer h.Hub.RemoveClient(conn)

		select {}
	})(c)
}

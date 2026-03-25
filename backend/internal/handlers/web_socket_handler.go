package handlers

import (
	"log"

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

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("WebSocet disconnected: ", err)
			}
			log.Println("Recieved from client: ", string(msg))
		}
	})(c)
}

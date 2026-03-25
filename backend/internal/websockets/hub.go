package websockets

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
	"github.com/magwach/distributed-task-scheduler/backend/internal/models"
)

type Hub struct {
	clients map[*websocket.Conn]bool
	lock    sync.Mutex
}

func HubInit() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (h *Hub) AddClient(conn *websocket.Conn) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.clients[conn] = true
}

func (h *Hub) RemoveClient(conn *websocket.Conn) {
	h.lock.Lock()
	defer h.lock.Unlock()

	delete(h.clients, conn)

	conn.Close()
}

func (h *Hub) Broadcast(message models.TaskUpdateEvent) {
	h.lock.Lock()
	defer h.lock.Unlock()

	data, err := json.Marshal(message)

	if err != nil {
		log.Println("Failed to parse the message to JSON")
		return
	}

	for conn := range h.clients {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			delete(h.clients, conn)
			conn.Close()
		}
	}
}

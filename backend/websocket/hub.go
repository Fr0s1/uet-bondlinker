package websocket

import (
	"encoding/json"
	"sync"

	"github.com/google/uuid"
)

// Client represents a connected websocket client
type Client struct {
	Hub      *Hub
	UserID   uuid.UUID
	Send     chan []byte
	IsTyping bool
}

// Hub maintains active clients and broadcasts messages
type Hub struct {
	clients    map[*Client]bool
	userConns  map[uuid.UUID][]*Client // Map user ID to their active connections
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
}

// Event represents a websocket event
type Event struct {
	ToUserId uuid.UUID       `json:"toUserId"`
	Type     string          `json:"type"`
	Payload  json.RawMessage `json:"payload"`
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		userConns:  make(map[uuid.UUID][]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.userConns[client.UserID] = append(h.userConns[client.UserID], client)
			h.mutex.Unlock()

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				h.removeUserConnection(client)
				close(client.Send)
			}
			h.mutex.Unlock()

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					h.mutex.Lock()
					close(client.Send)
					delete(h.clients, client)
					h.removeUserConnection(client)
					h.mutex.Unlock()
				}
			}
		}
	}
}

// SendToUser sends a message to all connections of a specific user
func (h *Hub) SendToUser(userID uuid.UUID, message []byte) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if clients, exists := h.userConns[userID]; exists {
		for _, client := range clients {
			select {
			case client.Send <- message:
			default:
				// If send fails, connection will be cleaned up in broadcast loop
			}
		}
	}
}

func (h *Hub) removeUserConnection(client *Client) {
	if conns, exists := h.userConns[client.UserID]; exists {
		var newConns []*Client
		for _, conn := range conns {
			if conn != client {
				newConns = append(newConns, conn)
			}
		}
		if len(newConns) == 0 {
			delete(h.userConns, client.UserID)
		} else {
			h.userConns[client.UserID] = newConns
		}
	}
}

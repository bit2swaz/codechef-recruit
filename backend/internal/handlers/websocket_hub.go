package handlers

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	RoomID   string
	PlayerID string
	Send     chan []byte
}

type Hub struct {
	rooms      map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *BroadcastMessage
	mu         sync.RWMutex
}

type BroadcastMessage struct {
	RoomID  string
	Message []byte
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *BroadcastMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.rooms[client.RoomID] == nil {
				h.rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.rooms[client.RoomID][client] = true
			h.mu.Unlock()
			log.Printf("Client registered to room %s (PlayerID: %s). Total clients in room: %d",
				client.RoomID, client.PlayerID, len(h.rooms[client.RoomID]))

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.rooms[client.RoomID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.Send)
					log.Printf("Client unregistered from room %s (PlayerID: %s). Remaining clients: %d",
						client.RoomID, client.PlayerID, len(clients))

					// Clean up empty rooms
					if len(clients) == 0 {
						delete(h.rooms, client.RoomID)
						log.Printf("Room %s is now empty and removed", client.RoomID)
					}
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			clients := h.rooms[message.RoomID]
			clientList := make([]*Client, 0, len(clients))
			for client := range clients {
				clientList = append(clientList, client)
			}
			h.mu.RUnlock()

			for _, client := range clientList {
				select {
				case client.Send <- message.Message:
				default:
					close(client.Send)
					h.mu.Lock()
					if roomClients, ok := h.rooms[message.RoomID]; ok {
						delete(roomClients, client)
					}
					h.mu.Unlock()
				}
			}
		}
	}
}

func (h *Hub) BroadcastToRoom(roomID string, message []byte) {
	h.broadcast <- &BroadcastMessage{
		RoomID:  roomID,
		Message: message,
	}
}

func (h *Hub) GetClientCount(roomID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.rooms[roomID]; ok {
		return len(clients)
	}
	return 0
}

func (h *Hub) GetAllRoomIDs() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	roomIDs := make([]string, 0, len(h.rooms))
	for roomID := range h.rooms {
		roomIDs = append(roomIDs, roomID)
	}
	return roomIDs
}

var hub = NewHub()

func GetHub() *Hub {
	return hub
}

func InitHub() {
	go hub.Run()
	log.Println("WebSocket Hub initialized and running")
}

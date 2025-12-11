package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type WSMessage struct {
	Type      string                 `json:"type"`
	PlayerID  string                 `json:"playerId"`
	RoomID    string                 `json:"roomId"`
	Data      map[string]interface{} `json:"data"`
	Timestamp int64                  `json:"timestamp"`
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["roomId"]

	if roomID == "" {
		http.Error(w, "roomId is required", http.StatusBadRequest)
		return
	}

	playerID := r.URL.Query().Get("playerId")
	if playerID == "" {
		http.Error(w, "playerId query parameter is required", http.StatusBadRequest)
		return
	}

	room := roomManager.GetRoom(roomID)
	if room == nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		Conn:     conn,
		RoomID:   roomID,
		PlayerID: playerID,
		Send:     make(chan []byte, 256),
	}

	hub.register <- client

	welcomeMsg := WSMessage{
		Type:      "connected",
		PlayerID:  playerID,
		RoomID:    roomID,
		Data:      map[string]interface{}{"message": "Connected to room"},
		Timestamp: time.Now().Unix(),
	}
	welcomeJSON, _ := json.Marshal(welcomeMsg)
	client.Send <- welcomeJSON

	joinMsg := WSMessage{
		Type:      "player_joined",
		PlayerID:  playerID,
		RoomID:    roomID,
		Data:      map[string]interface{}{"playerCount": hub.GetClientCount(roomID)},
		Timestamp: time.Now().Unix(),
	}
	joinJSON, _ := json.Marshal(joinMsg)
	hub.BroadcastToRoom(roomID, joinJSON)

	go client.writePump()
	go client.readPump(hub)
}

func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.Conn.Close()

		leaveMsg := WSMessage{
			Type:      "player_left",
			PlayerID:  c.PlayerID,
			RoomID:    c.RoomID,
			Data:      map[string]interface{}{"playerCount": hub.GetClientCount(c.RoomID)},
			Timestamp: time.Now().Unix(),
		}
		leaveJSON, _ := json.Marshal(leaveMsg)
		hub.BroadcastToRoom(c.RoomID, leaveJSON)
	}()

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	c.Conn.SetReadLimit(maxMessageSize)

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var wsMsg WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			log.Printf("Error parsing WebSocket message: %v", err)
			continue
		}

		wsMsg.PlayerID = c.PlayerID
		wsMsg.RoomID = c.RoomID
		wsMsg.Timestamp = time.Now().Unix()

		msgJSON, _ := json.Marshal(wsMsg)
		hub.BroadcastToRoom(c.RoomID, msgJSON)

		log.Printf("Message from %s in room %s: %s", c.PlayerID, c.RoomID, wsMsg.Type)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

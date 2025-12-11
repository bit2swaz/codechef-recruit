package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type WSMessage struct {
	Type      string                 `json:"type"`
	PlayerID  string                 `json:"playerId"`
	RoomID    string                 `json:"roomId"`
	Data      map[string]interface{} `json:"data"`
	Timestamp int64                  `json:"timestamp"`
}

func main() {
	// Parse command line flags
	roomID := flag.String("room", "TEST", "Room ID to join")
	playerID := flag.String("player", "player-1", "Player ID")
	serverAddr := flag.String("addr", "localhost:8080", "Server address")
	flag.Parse()

	// Setup interrupt handler
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Build WebSocket URL
	u := url.URL{Scheme: "ws", Host: *serverAddr, Path: "/ws/" + *roomID}
	q := u.Query()
	q.Set("playerId", *playerID)
	u.RawQuery = q.Encode()

	log.Printf("Connecting to %s", u.String())

	// Connect to WebSocket
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// Channel for messages to send
	done := make(chan struct{})

	// Read messages from server
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			var wsMsg WSMessage
			if err := json.Unmarshal(message, &wsMsg); err != nil {
				log.Printf("Received (raw): %s", message)
			} else {
				log.Printf("Received [%s]: type=%s, playerID=%s, data=%v",
					time.Unix(wsMsg.Timestamp, 0).Format("15:04:05"),
					wsMsg.Type,
					wsMsg.PlayerID,
					wsMsg.Data)
			}
		}
	}()

	// Send a test message after connection
	time.Sleep(1 * time.Second)
	testMsg := WSMessage{
		Type: "test_message",
		Data: map[string]interface{}{
			"message": "Hello from " + *playerID,
		},
	}
	msgJSON, _ := json.Marshal(testMsg)
	err = c.WriteMessage(websocket.TextMessage, msgJSON)
	if err != nil {
		log.Println("write:", err)
		return
	}
	log.Println("Sent test message")

	// Wait for interrupt or connection close
	select {
	case <-done:
		log.Println("Connection closed")
	case <-interrupt:
		log.Println("Interrupt received, closing connection")

		// Send close message
		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("write close:", err)
			return
		}

		// Wait for server to close or timeout
		select {
		case <-done:
		case <-time.After(time.Second):
		}
	}
}

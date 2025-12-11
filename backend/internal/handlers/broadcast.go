package handlers

import (
	"encoding/json"
	"log"

	"github.com/bit2swaz/codechef-recruit/backend/internal/store"
)

type GameMessage struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

// Broadcast sends a message to all connected clients in a room
func Broadcast(roomID string, messageType string, payload map[string]interface{}) {
	message := GameMessage{
		Type:    messageType,
		Payload: payload,
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling broadcast message: %v", err)
		return
	}

	hub := GetHub()
	hub.BroadcastToRoom(roomID, messageJSON)
	log.Printf("Broadcast to room %s: type=%s, payload=%v", roomID, messageType, payload)
}

// SendToPlayer sends a message to a specific player in a room
func SendToPlayer(roomID string, playerID string, messageType string, payload map[string]interface{}) {
	message := GameMessage{
		Type:    messageType,
		Payload: payload,
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling player message: %v", err)
		return
	}

	hub := GetHub()

	hub.mu.RLock()
	clients := hub.rooms[roomID]
	hub.mu.RUnlock()

	for client := range clients {
		if client.PlayerID == playerID {
			select {
			case client.Send <- messageJSON:
				log.Printf("Sent to player %s in room %s: type=%s", playerID, roomID, messageType)
			default:
				log.Printf("Failed to send to player %s (channel full)", playerID)
			}
			break
		}
	}
}

func BroadcastPlayerJoined(roomID string, playerName string, playerID string) {
	Broadcast(roomID, "PLAYER_JOINED", map[string]interface{}{
		"name":     playerName,
		"playerId": playerID,
	})
}

func BroadcastGameStart(roomID string) {
	Broadcast(roomID, "GAME_START", map[string]interface{}{
		"message": "All players ready! Roles have been assigned.",
	})
}

func SendRoleToPlayer(roomID string, player store.Player) {
	SendToPlayer(roomID, player.ID, "YOUR_ROLE", map[string]interface{}{
		"role": player.Role,
		"name": player.Name,
	})
}

func BroadcastRolesAssigned(roomID string, players []store.Player) {
	BroadcastGameStart(roomID)

	for _, player := range players {
		SendRoleToPlayer(roomID, player)
	}
}

func BroadcastGuessResult(roomID string, mantriName string, correct bool, scores map[string]int) {
	Broadcast(roomID, "GUESS_RESULT", map[string]interface{}{
		"mantri":  mantriName,
		"correct": correct,
		"scores":  scores,
	})
}

func BroadcastGameEnd(roomID string, finalScores map[string]interface{}) {
	Broadcast(roomID, "GAME_END", map[string]interface{}{
		"message": "Game finished!",
		"scores":  finalScores,
	})
}

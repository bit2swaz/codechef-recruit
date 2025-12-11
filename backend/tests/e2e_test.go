package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bit2swaz/codechef-recruit/backend/internal/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Player represents a test player with HTTP and WebSocket clients
type Player struct {
	Name         string
	ID           string
	WSConn       *websocket.Conn
	Messages     []map[string]interface{}
	MessageMutex sync.Mutex
	Role         string
}

// TestE2EFullGame simulates a complete game from start to finish
func TestE2EFullGame(t *testing.T) {
	// Step 1: Setup test server
	server := setupTestServer(t)
	defer server.Close()

	baseURL := server.URL
	wsURL := "ws" + strings.TrimPrefix(baseURL, "http")

	// Step 2: Create room via HTTP
	t.Log("=== Step 1: Creating room ===")
	roomID := createRoom(t, baseURL, "Alice")
	t.Logf("✓ Room created: %s", roomID)

	// Step 3: Create 4 players
	players := []*Player{
		{Name: "Alice", ID: ""},
		{Name: "Bob", ID: ""},
		{Name: "Charlie", ID: ""},
		{Name: "Diana", ID: ""},
	}

	// Step 4: Join all players via HTTP first to get player IDs
	t.Log("=== Step 2: Players joining room ===")

	// Alice is already in the room (creator), get her ID
	roomDetails := getRoomDetails(t, baseURL, roomID)
	if len(roomDetails.Players) != 1 {
		t.Fatalf("Expected 1 player after room creation, got %d", len(roomDetails.Players))
	}
	players[0].ID = roomDetails.Players[0].ID
	t.Logf("✓ Alice (creator) ID: %s", players[0].ID)

	// Join remaining 3 players
	for i := 1; i < 4; i++ {
		playerID := joinRoom(t, baseURL, roomID, players[i].Name)
		players[i].ID = playerID
		t.Logf("✓ %s joined with ID: %s", players[i].Name, playerID)
	}

	// Verify all 4 players are in the room
	roomDetails = getRoomDetails(t, baseURL, roomID)
	if len(roomDetails.Players) != 4 {
		t.Fatalf("Expected 4 players in room, got %d", len(roomDetails.Players))
	}
	t.Log("✓ All 4 players joined successfully")

	// Step 5: Connect all players to WebSocket in goroutines
	t.Log("=== Step 3: Connecting WebSocket clients ===")
	var wg sync.WaitGroup
	errChan := make(chan error, 4)

	for i, player := range players {
		wg.Add(1)
		go func(p *Player, index int) {
			defer wg.Done()

			// Connect to WebSocket (URL-encode the player ID)
			playerIDEncoded := url.QueryEscape(p.ID)
			wsURLFormatted := fmt.Sprintf("%s/ws/%s?playerId=%s", wsURL, roomID, playerIDEncoded)
			dialer := websocket.Dialer{}
			conn, _, err := dialer.Dial(wsURLFormatted, nil)
			if err != nil {
				errChan <- fmt.Errorf("%s failed to connect: %v", p.Name, err)
				return
			}
			p.WSConn = conn

			// Start reading messages in background
			go func() {
				for {
					_, message, err := conn.ReadMessage()
					if err != nil {
						return
					}

					var msg map[string]interface{}
					if err := json.Unmarshal(message, &msg); err != nil {
						continue
					}

					p.MessageMutex.Lock()
					p.Messages = append(p.Messages, msg)

					// Extract role if this is YOUR_ROLE message
					if msgType, ok := msg["type"].(string); ok && msgType == "YOUR_ROLE" {
						if payload, ok := msg["payload"].(map[string]interface{}); ok {
							if role, ok := payload["role"].(string); ok {
								p.Role = role
							}
						}
					}
					p.MessageMutex.Unlock()
				}
			}()

			t.Logf("✓ %s connected to WebSocket", p.Name)
		}(player, i)
	}

	wg.Wait()
	close(errChan)

	// Check for connection errors
	for err := range errChan {
		t.Fatal(err)
	}

	// Wait for welcome messages
	time.Sleep(200 * time.Millisecond)
	t.Log("✓ All WebSocket clients connected")

	// Step 6: Start the game
	t.Log("=== Step 4: Starting game ===")
	startGame(t, baseURL, roomID)
	t.Log("✓ Game start request sent")

	// Wait longer for GAME_START and YOUR_ROLE messages to be processed
	time.Sleep(1 * time.Second)

	// Step 7: Verify all players received GAME_START
	t.Log("=== Step 5: Verifying GAME_START broadcast ===")
	for _, player := range players {
		player.MessageMutex.Lock()
		gameStartReceived := false
		roleReceived := false

		for _, msg := range player.Messages {
			if msgType, ok := msg["type"].(string); ok {
				if msgType == "GAME_START" {
					gameStartReceived = true
				}
				if msgType == "YOUR_ROLE" {
					roleReceived = true
				}
			}
		}
		player.MessageMutex.Unlock()

		if !gameStartReceived {
			t.Errorf("%s did not receive GAME_START message", player.Name)
		} else {
			t.Logf("✓ %s received GAME_START", player.Name)
		}

		if !roleReceived {
			t.Errorf("%s did not receive YOUR_ROLE message", player.Name)
		} else {
			t.Logf("✓ %s received YOUR_ROLE: %s", player.Name, player.Role)
		}
	}

	// Step 8: Find the Mantri and Chor
	t.Log("=== Step 6: Finding Mantri and Chor ===")
	var mantri, chor *Player
	for _, player := range players {
		if player.Role == "Mantri" {
			mantri = player
			t.Logf("✓ Mantri is: %s", player.Name)
		}
		if player.Role == "Chor" {
			chor = player
			t.Logf("✓ Chor is: %s", player.Name)
		}
	}

	if mantri == nil || chor == nil {
		t.Fatal("Could not find Mantri or Chor")
	}

	// Step 9: Mantri makes a guess
	t.Log("=== Step 7: Mantri making guess ===")

	// Mantri guesses the Chor (correct guess)
	guessResult := submitGuess(t, baseURL, roomID, mantri.ID, chor.ID)
	t.Logf("✓ Mantri guessed, result: correct=%v", guessResult.Correct)

	// Wait longer for GUESS_RESULT and GAME_END messages
	time.Sleep(1 * time.Second)

	// Step 10: Verify all players received GUESS_RESULT and GAME_END
	t.Log("=== Step 8: Verifying GUESS_RESULT and GAME_END broadcasts ===")
	for _, player := range players {
		player.MessageMutex.Lock()
		guessResultReceived := false
		gameEndReceived := false

		for _, msg := range player.Messages {
			if msgType, ok := msg["type"].(string); ok {
				if msgType == "GUESS_RESULT" {
					guessResultReceived = true
					if payload, ok := msg["payload"].(map[string]interface{}); ok {
						t.Logf("  %s received GUESS_RESULT: correct=%v", player.Name, payload["correct"])
					}
				}
				if msgType == "GAME_END" {
					gameEndReceived = true
					if payload, ok := msg["payload"].(map[string]interface{}); ok {
						if scores, ok := payload["scores"].(map[string]interface{}); ok {
							t.Logf("  %s received GAME_END with scores: %v", player.Name, scores)
						}
					}
				}
			}
		}
		player.MessageMutex.Unlock()

		if !guessResultReceived {
			t.Errorf("%s did not receive GUESS_RESULT message", player.Name)
		} else {
			t.Logf("✓ %s received GUESS_RESULT", player.Name)
		}

		if !gameEndReceived {
			t.Errorf("%s did not receive GAME_END message", player.Name)
		} else {
			t.Logf("✓ %s received GAME_END", player.Name)
		}
	}

	// Step 11: Verify final scores
	t.Log("=== Step 9: Verifying final scores ===")
	if guessResult.Correct {
		t.Log("Correct guess - verifying scoring:")
		// Raja should have 1000, Mantri 800, Sipahi 500, Chor 0
		for _, player := range players {
			score := guessResult.UpdatedScores[player.ID]
			expectedScore := 0
			switch player.Role {
			case "Raja":
				expectedScore = 1000
			case "Mantri":
				expectedScore = 800
			case "Sipahi":
				expectedScore = 500
			case "Chor":
				expectedScore = 0
			}

			if score != expectedScore {
				t.Errorf("%s (%s) expected score %d, got %d", player.Name, player.Role, expectedScore, score)
			} else {
				t.Logf("✓ %s (%s): %d points", player.Name, player.Role, score)
			}
		}
	}

	// Step 12: Verify room status is FINISHED
	t.Log("=== Step 10: Verifying room status ===")
	finalRoomDetails := getRoomDetails(t, baseURL, roomID)
	if finalRoomDetails.Status != "FINISHED" {
		t.Errorf("Expected room status FINISHED, got %s", finalRoomDetails.Status)
	} else {
		t.Logf("✓ Room status: %s", finalRoomDetails.Status)
	}

	// Cleanup: Close all WebSocket connections
	t.Log("=== Cleanup: Closing connections ===")
	for _, player := range players {
		if player.WSConn != nil {
			player.WSConn.Close()
		}
	}

	t.Log("=== ✅ E2E Test Complete ===")
}

// Helper functions

func setupTestServer(t *testing.T) *httptest.Server {
	handlers.InitHub()

	r := mux.NewRouter()
	r.HandleFunc("/room/create", handlers.CreateRoom).Methods("POST")
	r.HandleFunc("/room/join", handlers.JoinRoom).Methods("POST")
	r.HandleFunc("/room/{roomId}", handlers.GetRoom).Methods("GET")
	r.HandleFunc("/game/start", handlers.StartGame).Methods("POST")
	r.HandleFunc("/game/guess", handlers.SubmitGuess).Methods("POST")
	r.HandleFunc("/ws/{roomId}", handlers.HandleWebSocket).Methods("GET")

	return httptest.NewServer(r)
}

func createRoom(t *testing.T, baseURL, playerName string) string {
	payload := map[string]string{"playerName": playerName}
	jsonData, _ := json.Marshal(payload)

	resp, err := http.Post(baseURL+"/room/create", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create room: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d", resp.StatusCode)
	}

	var result struct {
		RoomID string `json:"roomId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	return result.RoomID
}

func joinRoom(t *testing.T, baseURL, roomID, playerName string) string {
	payload := map[string]string{
		"roomId":     roomID,
		"playerName": playerName,
	}
	jsonData, _ := json.Marshal(payload)

	resp, err := http.Post(baseURL+"/room/join", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to join room: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, string(body))
	}

	// Get the updated room details to find the new player's ID
	time.Sleep(100 * time.Millisecond) // Give server time to process
	roomDetails := getRoomDetails(t, baseURL, roomID)

	// Find the player by name
	for _, player := range roomDetails.Players {
		if player.Name == playerName {
			return player.ID
		}
	}

	t.Fatalf("Could not find player %s in room after joining", playerName)
	return ""
}

type RoomDetails struct {
	RoomID  string `json:"roomId"`
	Status  string `json:"status"`
	Players []struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Score int    `json:"score"`
	} `json:"players"`
}

func getRoomDetails(t *testing.T, baseURL, roomID string) RoomDetails {
	resp, err := http.Get(baseURL + "/room/" + roomID)
	if err != nil {
		t.Fatalf("Failed to get room details: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var details RoomDetails
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		t.Fatalf("Failed to decode room details: %v", err)
	}

	return details
}

func startGame(t *testing.T, baseURL, roomID string) {
	payload := map[string]string{"roomId": roomID}
	jsonData, _ := json.Marshal(payload)

	resp, err := http.Post(baseURL+"/game/start", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, string(body))
	}
}

type GuessResult struct {
	Correct       bool           `json:"correct"`
	MantriID      string         `json:"mantriId"`
	ChorID        string         `json:"chorId"`
	ActualChorID  string         `json:"actualChorId"`
	UpdatedScores map[string]int `json:"updatedScores"`
}

func submitGuess(t *testing.T, baseURL, roomID, mantriID, guessedChorID string) GuessResult {
	payload := map[string]string{
		"roomId":              roomID,
		"mantriPlayerId":      mantriID,
		"guessedChorPlayerId": guessedChorID,
	}
	jsonData, _ := json.Marshal(payload)

	resp, err := http.Post(baseURL+"/game/guess", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to submit guess: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, string(body))
	}

	var result GuessResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode guess result: %v", err)
	}

	return result
}

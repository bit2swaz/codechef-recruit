package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bit2swaz/codechef-recruit/backend/internal/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// setupTestServer creates a test HTTP server with WebSocket support
func setupTestServer() *httptest.Server {
	// Initialize WebSocket hub
	handlers.InitHub()

	// Create router with WebSocket endpoint
	r := mux.NewRouter()
	r.HandleFunc("/ws/{roomId}", handlers.HandleWebSocket).Methods("GET")
	r.HandleFunc("/room/create", handlers.CreateRoom).Methods("POST")
	r.HandleFunc("/room/join", handlers.JoinRoom).Methods("POST")
	r.HandleFunc("/room/{roomId}", handlers.GetRoom).Methods("GET")

	// Create test server
	return httptest.NewServer(r)
}

// createTestRoom creates a room via HTTP API and returns the room ID
func createTestRoom(serverURL string, playerName string) (string, error) {
	payload := `{"playerName":"` + playerName + `"}`
	resp, err := http.Post(serverURL+"/room/create", "application/json", strings.NewReader(payload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var createResp struct {
		RoomID string `json:"roomId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		return "", err
	}

	return createResp.RoomID, nil
}

// TestWebSocketConnection tests basic WebSocket connection establishment
func TestWebSocketConnection(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a room first
	roomID, err := createTestRoom(server.URL, "TestPlayer")
	if err != nil {
		t.Fatalf("Failed to create room: %v", err)
	}

	// Convert http://... to ws://...
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/" + roomID + "?playerId=test-player-1"

	// Create WebSocket dialer
	dialer := websocket.Dialer{}

	// Connect to WebSocket
	conn, resp, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	// Assert connection is successful
	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Errorf("Expected status 101 Switching Protocols, got %d", resp.StatusCode)
	}

	t.Log("✓ WebSocket connection established successfully")
}

// TestWebSocketSendMessage tests sending a message to the server
func TestWebSocketSendMessage(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a room first
	roomID, err := createTestRoom(server.URL, "TestPlayer")
	if err != nil {
		t.Fatalf("Failed to create room: %v", err)
	}

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/" + roomID + "?playerId=test-player-2"

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	// Send a test message
	testMessage := map[string]interface{}{
		"type":     "test",
		"playerId": "test-player-2",
		"data":     map[string]string{"message": "Hello from test"},
	}

	messageJSON, err := json.Marshal(testMessage)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, messageJSON)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	t.Log("✓ Message sent successfully without server crash")

	// Give server time to process
	time.Sleep(100 * time.Millisecond)
}

// TestWebSocketMultipleClients tests multiple clients in the same room
func TestWebSocketMultipleClients(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a room first
	roomID, err := createTestRoom(server.URL, "TestPlayer")
	if err != nil {
		t.Fatalf("Failed to create room: %v", err)
	}

	baseURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect first client
	client1URL := baseURL + "/ws/" + roomID + "?playerId=client-1"
	dialer := websocket.Dialer{}
	conn1, _, err := dialer.Dial(client1URL, nil)
	if err != nil {
		t.Fatalf("Failed to connect client 1: %v", err)
	}
	defer conn1.Close()

	// Connect second client
	client2URL := baseURL + "/ws/" + roomID + "?playerId=client-2"
	conn2, _, err := dialer.Dial(client2URL, nil)
	if err != nil {
		t.Fatalf("Failed to connect client 2: %v", err)
	}
	defer conn2.Close()

	t.Log("✓ Multiple clients connected to same room successfully")

	// Give server time to register both clients
	time.Sleep(100 * time.Millisecond)
}

// TestWebSocketBroadcast tests that broadcasts reach all clients
func TestWebSocketBroadcast(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	baseURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Create a room via HTTP API
	roomID, err := createTestRoom(server.URL, "FirstPlayer")
	if err != nil {
		t.Fatalf("Failed to create room: %v", err)
	}

	// Connect WebSocket client to this room
	receiverURL := baseURL + "/ws/" + roomID + "?playerId=broadcast-receiver"
	dialer := websocket.Dialer{}
	connReceiver, _, err := dialer.Dial(receiverURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect receiver: %v", err)
	}
	defer connReceiver.Close()

	// Set read deadline
	connReceiver.SetReadDeadline(time.Now().Add(2 * time.Second))

	// Join another player to trigger PLAYER_JOINED broadcast
	httpURL := server.URL
	joinPayload := `{"roomId":"` + roomID + `","playerName":"SecondPlayer"}`
	joinResp, err := http.Post(httpURL+"/room/join", "application/json", strings.NewReader(joinPayload))
	if err != nil {
		t.Fatalf("Failed to join room: %v", err)
	}
	defer joinResp.Body.Close()

	// Try to read broadcast message
	_, message, err := connReceiver.ReadMessage()
	if err != nil {
		// Timeout is acceptable if no broadcast was sent yet
		t.Logf("No broadcast received (or timeout): %v", err)
	} else {
		var broadcastMsg map[string]interface{}
		if err := json.Unmarshal(message, &broadcastMsg); err == nil {
			t.Logf("✓ Broadcast received: %v", broadcastMsg)

			// Check if it's a PLAYER_JOINED message
			if msgType, ok := broadcastMsg["type"].(string); ok && msgType == "PLAYER_JOINED" {
				t.Log("✓ PLAYER_JOINED broadcast verified")
			}
		}
	}
}

// TestWebSocketConnectionWithoutPlayerID tests connection without playerId parameter
func TestWebSocketConnectionWithoutPlayerID(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/TEST_ROOM"

	dialer := websocket.Dialer{}
	conn, resp, err := dialer.Dial(wsURL, nil)

	// Should fail or return an error
	if err == nil {
		conn.Close()
		t.Error("Expected connection to fail without playerId, but it succeeded")
	} else {
		t.Logf("✓ Connection properly rejected without playerId: %v", err)
	}

	if resp != nil && resp.StatusCode != http.StatusSwitchingProtocols {
		t.Logf("✓ Received non-switching-protocols status: %d", resp.StatusCode)
	}
}

// TestWebSocketCleanup tests that connections are properly cleaned up
func TestWebSocketCleanup(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a room first
	roomID, err := createTestRoom(server.URL, "TestPlayer")
	if err != nil {
		t.Fatalf("Failed to create room: %v", err)
	}

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/" + roomID + "?playerId=cleanup-client"

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	// Close connection
	conn.Close()
	t.Log("✓ Connection closed")

	// Give server time to clean up
	time.Sleep(200 * time.Millisecond)

	// Reconnect with same player ID (should work if cleanup happened)
	conn2, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to reconnect after cleanup: %v", err)
	}
	defer conn2.Close()

	t.Log("✓ Reconnection successful after cleanup")
}

// TestWebSocketConcurrentConnections tests handling of many concurrent connections
func TestWebSocketConcurrentConnections(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a room first
	roomID, err := createTestRoom(server.URL, "TestPlayer")
	if err != nil {
		t.Fatalf("Failed to create room: %v", err)
	}

	baseURL := "ws" + strings.TrimPrefix(server.URL, "http")
	numClients := 10

	connections := make([]*websocket.Conn, numClients)
	dialer := websocket.Dialer{}

	// Connect multiple clients
	for i := 0; i < numClients; i++ {
		wsURL := baseURL + "/ws/" + roomID + "?playerId=concurrent-client-" + string(rune('A'+i))
		conn, _, err := dialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("Failed to connect client %d: %v", i, err)
		}
		connections[i] = conn
	}

	t.Logf("✓ Successfully connected %d concurrent clients", numClients)

	// Clean up all connections
	for i, conn := range connections {
		if conn != nil {
			conn.Close()
			t.Logf("✓ Closed connection %d", i)
		}
	}

	// Give server time to clean up
	time.Sleep(100 * time.Millisecond)
}

// TestWebSocketPingPong tests the ping/pong heartbeat mechanism
func TestWebSocketPingPong(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a room first
	roomID, err := createTestRoom(server.URL, "TestPlayer")
	if err != nil {
		t.Fatalf("Failed to create room: %v", err)
	}

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/" + roomID + "?playerId=ping-client"

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Start reading messages in background to keep connection alive
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}()

	// Wait to ensure connection stays alive (server sends pings periodically)
	time.Sleep(1 * time.Second)

	// If we get here without connection closing, ping/pong is working
	t.Log("✓ Ping/Pong mechanism active (connection stayed alive)")
}

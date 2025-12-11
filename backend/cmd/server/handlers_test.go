package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bit2swaz/codechef-recruit/backend/internal/handlers"
	"github.com/gorilla/mux"
)

func setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/room/create", handlers.CreateRoom).Methods("POST")
	r.HandleFunc("/room/join", handlers.JoinRoom).Methods("POST")
	r.HandleFunc("/room/{roomId}", handlers.GetRoom).Methods("GET")
	return r
}

// TestCreateRoom tests POST /room/create endpoint
func TestCreateRoom(t *testing.T) {
	router := setupRouter()

	// Create request body
	reqBody := map[string]string{
		"playerName": "Alice",
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "/room/create", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Check response body
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Verify roomId exists in response
	roomID, exists := response["roomId"]
	if !exists {
		t.Error("Response does not contain roomId")
	}

	// Verify roomId is not empty
	if roomID == "" {
		t.Error("roomId is empty")
	}

	// Verify roomId is 4 characters (as per requirement)
	if len(roomID) != 4 {
		t.Errorf("Expected roomId to be 4 characters, got %d", len(roomID))
	}

	t.Logf("Successfully created room with ID: %s", roomID)
}

// TestCreateRoomMissingPlayerName tests POST /room/create with missing playerName
func TestCreateRoomMissingPlayerName(t *testing.T) {
	router := setupRouter()

	// Create request body without playerName
	reqBody := map[string]string{}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/room/create", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Should return 400 Bad Request
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

// TestJoinRoom tests POST /room/join endpoint with success scenario
func TestJoinRoom(t *testing.T) {
	router := setupRouter()

	// Step 1: Create a room first
	createReqBody := map[string]string{
		"playerName": "Alice",
	}
	createJsonBody, err := json.Marshal(createReqBody)
	if err != nil {
		t.Fatalf("Failed to marshal create request body: %v", err)
	}

	createReq, err := http.NewRequest("POST", "/room/create", bytes.NewBuffer(createJsonBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	createReq.Header.Set("Content-Type", "application/json")

	createRR := httptest.NewRecorder()
	router.ServeHTTP(createRR, createReq)

	// Extract roomId from create response
	var createResponse map[string]string
	err = json.Unmarshal(createRR.Body.Bytes(), &createResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal create response: %v", err)
	}

	roomID := createResponse["roomId"]
	if roomID == "" {
		t.Fatal("Failed to get roomId from create response")
	}

	t.Logf("Created room with ID: %s", roomID)

	// Step 2: Join the created room
	joinReqBody := map[string]string{
		"roomId":     roomID,
		"playerName": "Bob",
	}
	joinJsonBody, err := json.Marshal(joinReqBody)
	if err != nil {
		t.Fatalf("Failed to marshal join request body: %v", err)
	}

	joinReq, err := http.NewRequest("POST", "/room/join", bytes.NewBuffer(joinJsonBody))
	if err != nil {
		t.Fatalf("Failed to create join request: %v", err)
	}
	joinReq.Header.Set("Content-Type", "application/json")

	joinRR := httptest.NewRecorder()
	router.ServeHTTP(joinRR, joinReq)

	// Check status code is 200 OK
	if status := joinRR.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check response body
	var joinResponse map[string]string
	err = json.Unmarshal(joinRR.Body.Bytes(), &joinResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal join response: %v", err)
	}

	// Verify response contains expected fields
	if joinResponse["roomId"] != roomID {
		t.Errorf("Expected roomId %s, got %s", roomID, joinResponse["roomId"])
	}

	if joinResponse["message"] == "" {
		t.Error("Expected message in response")
	}

	t.Logf("Successfully joined room %s with message: %s", roomID, joinResponse["message"])
}

// TestJoinRoomNonExistent tests POST /room/join with non-existent room (404)
func TestJoinRoomNonExistent(t *testing.T) {
	router := setupRouter()

	// Try to join a non-existent room
	joinReqBody := map[string]string{
		"roomId":     "XXXX",
		"playerName": "Bob",
	}
	jsonBody, err := json.Marshal(joinReqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/room/join", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check status code is 404 Not Found
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	// Check response contains error message
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if response["error"] == "" {
		t.Error("Expected error message in response")
	}

	t.Logf("Correctly returned 404 with error: %s", response["error"])
}

// TestJoinRoomFull tests POST /room/join when room is full (4 players max)
func TestJoinRoomFull(t *testing.T) {
	router := setupRouter()

	// Step 1: Create a room
	createReqBody := map[string]string{
		"playerName": "Player1",
	}
	createJsonBody, _ := json.Marshal(createReqBody)
	createReq, _ := http.NewRequest("POST", "/room/create", bytes.NewBuffer(createJsonBody))
	createReq.Header.Set("Content-Type", "application/json")

	createRR := httptest.NewRecorder()
	router.ServeHTTP(createRR, createReq)

	var createResponse map[string]string
	json.Unmarshal(createRR.Body.Bytes(), &createResponse)
	roomID := createResponse["roomId"]

	// Step 2: Add 3 more players (total 4 players including creator)
	for i := 2; i <= 4; i++ {
		joinReqBody := map[string]string{
			"roomId":     roomID,
			"playerName": "Player" + string(rune(i+48)),
		}
		joinJsonBody, _ := json.Marshal(joinReqBody)
		joinReq, _ := http.NewRequest("POST", "/room/join", bytes.NewBuffer(joinJsonBody))
		joinReq.Header.Set("Content-Type", "application/json")

		joinRR := httptest.NewRecorder()
		router.ServeHTTP(joinRR, joinReq)

		if joinRR.Code != http.StatusOK {
			t.Fatalf("Failed to add player %d", i)
		}
	}

	// Step 3: Try to add 5th player (should fail)
	joinReqBody := map[string]string{
		"roomId":     roomID,
		"playerName": "Player5",
	}
	joinJsonBody, _ := json.Marshal(joinReqBody)
	joinReq, _ := http.NewRequest("POST", "/room/join", bytes.NewBuffer(joinJsonBody))
	joinReq.Header.Set("Content-Type", "application/json")

	joinRR := httptest.NewRecorder()
	router.ServeHTTP(joinRR, joinReq)

	// Should return 400 Bad Request
	if status := joinRR.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	var response map[string]string
	json.Unmarshal(joinRR.Body.Bytes(), &response)

	if response["error"] == "" {
		t.Error("Expected error message when room is full")
	}

	t.Logf("Correctly rejected 5th player with error: %s", response["error"])
}

// TestGetRoom tests GET /room/{roomId} endpoint
func TestGetRoom(t *testing.T) {
	router := setupRouter()

	// Create a room first
	createReqBody := map[string]string{
		"playerName": "Alice",
	}
	createJsonBody, _ := json.Marshal(createReqBody)
	createReq, _ := http.NewRequest("POST", "/room/create", bytes.NewBuffer(createJsonBody))
	createReq.Header.Set("Content-Type", "application/json")

	createRR := httptest.NewRecorder()
	router.ServeHTTP(createRR, createReq)

	var createResponse map[string]string
	json.Unmarshal(createRR.Body.Bytes(), &createResponse)
	roomID := createResponse["roomId"]

	// Get the room details
	getReq, err := http.NewRequest("GET", "/room/"+roomID, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	getRR := httptest.NewRecorder()
	router.ServeHTTP(getRR, getReq)

	// Check status code
	if status := getRR.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check response structure
	var response map[string]interface{}
	err = json.Unmarshal(getRR.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify roomId
	if response["roomId"] != roomID {
		t.Errorf("Expected roomId %s, got %v", roomID, response["roomId"])
	}

	// Verify players array exists
	players, ok := response["players"].([]interface{})
	if !ok {
		t.Fatal("Expected players to be an array")
	}

	// Verify at least 1 player (the creator)
	if len(players) < 1 {
		t.Error("Expected at least 1 player in the room")
	}

	// Verify role is NOT exposed in player info
	firstPlayer, ok := players[0].(map[string]interface{})
	if !ok {
		t.Fatal("Expected player to be an object")
	}

	if _, hasRole := firstPlayer["role"]; hasRole {
		t.Error("Player info should NOT contain 'role' field (roles should be hidden)")
	}

	// Verify required fields are present
	if firstPlayer["name"] == "" {
		t.Error("Player should have a name")
	}

	t.Logf("Successfully retrieved room %s with %d players", roomID, len(players))
}

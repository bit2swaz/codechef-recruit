package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/bit2swaz/codechef-recruit/backend/internal/store"
	"github.com/gorilla/mux"
)

var (
	roomManager = store.NewRoomManager()
	random      = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type CreateRoomRequest struct {
	PlayerName string `json:"playerName"`
}

type CreateRoomResponse struct {
	RoomID string `json:"roomId"`
}

type JoinRoomRequest struct {
	RoomID     string `json:"roomId"`
	PlayerName string `json:"playerName"`
}

type JoinRoomResponse struct {
	Message string `json:"message"`
	RoomID  string `json:"roomId"`
}

type RoomDetailsResponse struct {
	RoomID  string             `json:"roomId"`
	Status  string             `json:"status"`
	Players []PlayerInfoPublic `json:"players"`
}

type PlayerInfoPublic struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func generateRoomID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 4)
	for i := range result {
		result[i] = charset[random.Intn(len(charset))]
	}
	return string(result)
}

func generatePlayerID() string {
	return time.Now().Format("20060102150405") + "-" + string(rune(random.Intn(10000)))
}

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	var req CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request body"})
		return
	}

	if req.PlayerName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "playerName is required"})
		return
	}

	roomID := generateRoomID()

	room := roomManager.CreateRoom(roomID)

	admin := store.Player{
		ID:    generatePlayerID(),
		Name:  req.PlayerName,
		Role:  "Admin",
		Score: 0,
	}
	room.AddPlayer(admin)

	BroadcastPlayerJoined(roomID, admin.Name, admin.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateRoomResponse{RoomID: roomID})
}

func JoinRoom(w http.ResponseWriter, r *http.Request) {
	var req JoinRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request body"})
		return
	}

	if req.RoomID == "" || req.PlayerName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "roomId and playerName are required"})
		return
	}

	room := roomManager.GetRoom(req.RoomID)
	if room == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Room not found"})
		return
	}

	players := room.GetPlayers()
	if len(players) >= 4 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Room is full (max 4 players)"})
		return
	}

	player := store.Player{
		ID:    generatePlayerID(),
		Name:  req.PlayerName,
		Role:  "Player",
		Score: 0,
	}
	room.AddPlayer(player)

	// Broadcast to all connected clients that a player joined
	BroadcastPlayerJoined(req.RoomID, player.Name, player.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JoinRoomResponse{
		Message: "Successfully joined room",
		RoomID:  req.RoomID,
	})
}

func GetRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["roomId"]

	if roomID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "roomId is required"})
		return
	}

	room := roomManager.GetRoom(roomID)
	if room == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Room not found"})
		return
	}

	players := room.GetPlayers()
	publicPlayers := make([]PlayerInfoPublic, len(players))
	for i, p := range players {
		publicPlayers[i] = PlayerInfoPublic{
			ID:    p.ID,
			Name:  p.Name,
			Score: p.Score,
		}
	}

	response := RoomDetailsResponse{
		RoomID:  room.ID,
		Status:  room.GetStatus(),
		Players: publicPlayers,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

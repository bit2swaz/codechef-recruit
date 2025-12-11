package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bit2swaz/codechef-recruit/backend/internal/game"
)

type StartGameRequest struct {
	RoomID string `json:"roomId"`
}

func StartGame(w http.ResponseWriter, r *http.Request) {
	var req StartGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request body"})
		return
	}

	if req.RoomID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "roomId is required"})
		return
	}

	room := roomManager.GetRoom(req.RoomID)
	if room == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Room not found"})
		return
	}

	players := room.GetPlayers()
	if len(players) != 4 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Need exactly 4 players to start game"})
		return
	}

	game.AssignRoles(room)

	players = room.GetPlayers()
	BroadcastRolesAssigned(room.ID, players)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Game started"})
}

type SubmitGuessRequest struct {
	RoomID              string `json:"roomId"`
	MantriPlayerID      string `json:"mantriPlayerId"`
	GuessedChorPlayerID string `json:"guessedChorPlayerId"`
}

func SubmitGuess(w http.ResponseWriter, r *http.Request) {
	var req SubmitGuessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request body"})
		return
	}

	if req.RoomID == "" || req.MantriPlayerID == "" || req.GuessedChorPlayerID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "roomId, mantriPlayerId, and guessedChorPlayerId are required"})
		return
	}

	room := roomManager.GetRoom(req.RoomID)
	if room == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Room not found"})
		return
	}

	result, err := game.ProcessGuess(room, req.MantriPlayerID, req.GuessedChorPlayerID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	players := room.GetPlayers()
	var mantriName string
	for _, p := range players {
		if p.ID == req.MantriPlayerID {
			mantriName = p.Name
			break
		}
	}
	BroadcastGuessResult(room.ID, mantriName, result.Correct, result.UpdatedScores)

	finalScores := make(map[string]interface{})
	for _, player := range players {
		finalScores[player.Name] = player.Score
	}
	BroadcastGameEnd(room.ID, finalScores)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

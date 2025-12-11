package game

import (
	"errors"
	"math/rand/v2"

	"github.com/bit2swaz/codechef-recruit/backend/internal/store"
)

func AssignRoles(room *store.Room) {
	players := room.GetPlayers()
	if len(players) != 4 {
		return
	}

	roles := []string{"Raja", "Mantri", "Chor", "Sipahi"}

	rand.Shuffle(len(roles), func(i, j int) {
		roles[i], roles[j] = roles[j], roles[i]
	})

	updatedPlayers := make([]store.Player, len(players))
	for i := 0; i < len(players); i++ {
		updatedPlayers[i] = players[i]
		updatedPlayers[i].Role = roles[i]
	}

	room.UpdatePlayersAndStatus(updatedPlayers, "GUESSING")
}

type GuessResult struct {
	Correct       bool
	MantriID      string
	ChorID        string
	ActualChorID  string
	UpdatedScores map[string]int
}

func ProcessGuess(room *store.Room, mantriPlayerID string, guessedChorPlayerID string) (*GuessResult, error) {
	players := room.GetPlayers()

	var mantri, chor, raja, sipahi *store.Player

	for i := range players {
		player := &players[i]
		switch player.Role {
		case "Mantri":
			mantri = player
		case "Chor":
			chor = player
		case "Raja":
			raja = player
		case "Sipahi":
			sipahi = player
		}
	}

	if mantri == nil || chor == nil || raja == nil || sipahi == nil {
		return nil, errors.New("room does not have all required roles assigned")
	}

	if mantri.ID != mantriPlayerID {
		return nil, errors.New("only the Mantri can make a guess")
	}

	correctGuess := (guessedChorPlayerID == chor.ID)

	updatedPlayers := make([]store.Player, len(players))
	copy(updatedPlayers, players)

	for i := range updatedPlayers {
		if updatedPlayers[i].Role == "Raja" {
			updatedPlayers[i].Score = 1000
		} else if updatedPlayers[i].Role == "Mantri" {
			if correctGuess {
				updatedPlayers[i].Score = 800
			} else {
				updatedPlayers[i].Score = 0
			}
		} else if updatedPlayers[i].Role == "Sipahi" {
			updatedPlayers[i].Score = 500
		} else if updatedPlayers[i].Role == "Chor" {
			if correctGuess {
				updatedPlayers[i].Score = 0
			} else {
				updatedPlayers[i].Score = 800
			}
		}
	}

	room.UpdatePlayersAndStatus(updatedPlayers, "FINISHED")

	result := &GuessResult{
		Correct:       correctGuess,
		MantriID:      mantri.ID,
		ChorID:        guessedChorPlayerID,
		ActualChorID:  chor.ID,
		UpdatedScores: make(map[string]int),
	}

	for _, player := range updatedPlayers {
		result.UpdatedScores[player.ID] = player.Score
	}

	return result, nil
}

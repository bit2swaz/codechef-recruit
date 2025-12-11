package game

import (
	"testing"

	"github.com/bit2swaz/codechef-recruit/backend/internal/store"
)

// TestScoringLogicTableDriven is a comprehensive table-driven test for scoring logic
func TestScoringLogicTableDriven(t *testing.T) {
	// Define test cases
	testCases := []struct {
		Name           string
		MantriID       string
		ChorID         string
		GuessedID      string
		ExpectedScores map[string]int // Role -> Score
		ShouldSucceed  bool
		ErrorContains  string
	}{
		{
			Name:      "Mantri Guesses Correctly",
			MantriID:  "mantri-correct",
			ChorID:    "chor-correct",
			GuessedID: "chor-correct",
			ExpectedScores: map[string]int{
				"Raja":   1000,
				"Mantri": 800,
				"Sipahi": 500,
				"Chor":   0,
			},
			ShouldSucceed: true,
		},
		{
			Name:      "Mantri Guesses Wrong - Guesses Sipahi",
			MantriID:  "mantri-wrong-1",
			ChorID:    "chor-wrong-1",
			GuessedID: "sipahi-wrong-1",
			ExpectedScores: map[string]int{
				"Raja":   1000,
				"Mantri": 0,
				"Sipahi": 500,
				"Chor":   800,
			},
			ShouldSucceed: true,
		},
		{
			Name:      "Mantri Guesses Wrong - Guesses Raja",
			MantriID:  "mantri-wrong-2",
			ChorID:    "chor-wrong-2",
			GuessedID: "raja-wrong-2",
			ExpectedScores: map[string]int{
				"Raja":   1000,
				"Mantri": 0,
				"Sipahi": 500,
				"Chor":   800,
			},
			ShouldSucceed: true,
		},
		{
			Name:      "Mantri Guesses Wrong - Guesses Self",
			MantriID:  "mantri-wrong-3",
			ChorID:    "chor-wrong-3",
			GuessedID: "mantri-wrong-3",
			ExpectedScores: map[string]int{
				"Raja":   1000,
				"Mantri": 0,
				"Sipahi": 500,
				"Chor":   800,
			},
			ShouldSucceed: true,
		},
		{
			Name:          "Non-Mantri Tries to Guess",
			MantriID:      "mantri-invalid",
			ChorID:        "chor-invalid",
			GuessedID:     "chor-invalid",
			ShouldSucceed: false,
			ErrorContains: "only the Mantri can make a guess",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Create a new room for each test case
			rm := store.NewRoomManager()
			room := rm.CreateRoom("TEST-" + tc.Name)

			// Set up players with specific roles
			var players []store.Player

			if tc.ShouldSucceed {
				// Valid scenario - create all 4 roles
				players = []store.Player{
					{ID: "raja-" + tc.Name, Name: "Raja", Role: "Raja", Score: 0},
					{ID: tc.MantriID, Name: "Mantri", Role: "Mantri", Score: 0},
					{ID: tc.ChorID, Name: "Chor", Role: "Chor", Score: 0},
					{ID: "sipahi-" + tc.Name, Name: "Sipahi", Role: "Sipahi", Score: 0},
				}
			} else {
				// Invalid scenario for error testing
				players = []store.Player{
					{ID: "raja-" + tc.Name, Name: "Raja", Role: "Raja", Score: 0},
					{ID: tc.MantriID, Name: "Mantri", Role: "Mantri", Score: 0},
					{ID: tc.ChorID, Name: "Chor", Role: "Chor", Score: 0},
					{ID: "sipahi-" + tc.Name, Name: "Sipahi", Role: "Sipahi", Score: 0},
				}
			}

			for _, player := range players {
				room.AddPlayer(player)
			}

			room.UpdatePlayersAndStatus(players, "GUESSING")

			// Determine who is making the guess
			var guesserID string
			if tc.ShouldSucceed {
				guesserID = tc.MantriID
			} else {
				// For error case, use a different player ID
				guesserID = "raja-" + tc.Name
			}

			// Process the guess
			result, err := ProcessGuess(room, guesserID, tc.GuessedID)

			// Check if error handling is correct
			if !tc.ShouldSucceed {
				if err == nil {
					t.Errorf("Expected error containing '%s', but got no error", tc.ErrorContains)
					return
				}
				if tc.ErrorContains != "" && err.Error() != tc.ErrorContains {
					t.Errorf("Expected error containing '%s', got '%s'", tc.ErrorContains, err.Error())
				}
				return
			}

			// For successful cases, verify no error
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result == nil {
				t.Fatal("Result is nil")
			}

			// Verify scores match expected values
			finalPlayers := room.GetPlayers()
			actualScoresByRole := make(map[string]int)

			for _, player := range finalPlayers {
				actualScoresByRole[player.Role] = player.Score
			}

			// Compare expected vs actual scores
			for role, expectedScore := range tc.ExpectedScores {
				actualScore, exists := actualScoresByRole[role]
				if !exists {
					t.Errorf("Role %s not found in actual scores", role)
					continue
				}
				if actualScore != expectedScore {
					t.Errorf("Role %s: expected score %d, got %d", role, expectedScore, actualScore)
				}
			}

			// Verify room status is FINISHED
			if room.GetStatus() != "FINISHED" {
				t.Errorf("Expected room status FINISHED, got %s", room.GetStatus())
			}

			// Verify total score is always 2300
			totalScore := 0
			for _, score := range actualScoresByRole {
				totalScore += score
			}
			if totalScore != 2300 {
				t.Errorf("Expected total score 2300, got %d", totalScore)
			}

			// Log the scores for visibility
			t.Logf("Scores - Raja: %d, Mantri: %d, Sipahi: %d, Chor: %d (Total: %d)",
				actualScoresByRole["Raja"],
				actualScoresByRole["Mantri"],
				actualScoresByRole["Sipahi"],
				actualScoresByRole["Chor"],
				totalScore)
		})
	}
}

// TestScoringEdgeCases tests edge cases in scoring
func TestScoringEdgeCases(t *testing.T) {
	testCases := []struct {
		Name          string
		SetupPlayers  func() []store.Player
		GuessFunc     func(players []store.Player) (mantriID, guessID string)
		ShouldSucceed bool
		ErrorContains string
	}{
		{
			Name: "Room with missing roles",
			SetupPlayers: func() []store.Player {
				return []store.Player{
					{ID: "player1", Name: "Player 1", Role: "", Score: 0},
					{ID: "player2", Name: "Player 2", Role: "", Score: 0},
					{ID: "player3", Name: "Player 3", Role: "", Score: 0},
					{ID: "player4", Name: "Player 4", Role: "", Score: 0},
				}
			},
			GuessFunc: func(players []store.Player) (string, string) {
				return "player1", "player2"
			},
			ShouldSucceed: false,
			ErrorContains: "room does not have all required roles assigned",
		},
		{
			Name: "Room with only 3 players",
			SetupPlayers: func() []store.Player {
				return []store.Player{
					{ID: "raja", Name: "Raja", Role: "Raja", Score: 0},
					{ID: "mantri", Name: "Mantri", Role: "Mantri", Score: 0},
					{ID: "chor", Name: "Chor", Role: "Chor", Score: 0},
				}
			},
			GuessFunc: func(players []store.Player) (string, string) {
				return "mantri", "chor"
			},
			ShouldSucceed: false,
			ErrorContains: "room does not have all required roles assigned",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			rm := store.NewRoomManager()
			room := rm.CreateRoom("EDGE-" + tc.Name)

			players := tc.SetupPlayers()
			for _, player := range players {
				room.AddPlayer(player)
			}

			room.UpdatePlayersAndStatus(players, "GUESSING")

			mantriID, guessID := tc.GuessFunc(players)
			result, err := ProcessGuess(room, mantriID, guessID)

			if !tc.ShouldSucceed {
				if err == nil {
					t.Errorf("Expected error, got nil")
					return
				}
				if tc.ErrorContains != "" && err.Error() != tc.ErrorContains {
					t.Errorf("Expected error '%s', got '%s'", tc.ErrorContains, err.Error())
				}
				if result != nil {
					t.Error("Expected nil result on error")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestScoringConsistency verifies scoring is consistent across multiple runs
func TestScoringConsistency(t *testing.T) {
	scenarios := []struct {
		name         string
		correctGuess bool
	}{
		{"correct_guess", true},
		{"wrong_guess", false},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// Run the same scenario multiple times
			for i := 0; i < 10; i++ {
				rm := store.NewRoomManager()
				room := rm.CreateRoom("CONSISTENCY")

				players := []store.Player{
					{ID: "raja", Name: "Raja", Role: "Raja", Score: 0},
					{ID: "mantri", Name: "Mantri", Role: "Mantri", Score: 0},
					{ID: "chor", Name: "Chor", Role: "Chor", Score: 0},
					{ID: "sipahi", Name: "Sipahi", Role: "Sipahi", Score: 0},
				}

				for _, player := range players {
					room.AddPlayer(player)
				}

				room.UpdatePlayersAndStatus(players, "GUESSING")

				var guessID string
				if scenario.correctGuess {
					guessID = "chor"
				} else {
					guessID = "sipahi"
				}

				result, err := ProcessGuess(room, "mantri", guessID)
				if err != nil {
					t.Fatalf("Iteration %d: unexpected error: %v", i, err)
				}

				// Verify correct/incorrect guess
				if result.Correct != scenario.correctGuess {
					t.Errorf("Iteration %d: expected correct=%v, got %v", i, scenario.correctGuess, result.Correct)
				}

				// Verify total is always 2300
				total := 0
				for _, score := range result.UpdatedScores {
					total += score
				}
				if total != 2300 {
					t.Errorf("Iteration %d: expected total 2300, got %d", i, total)
				}
			}
		})
	}
}

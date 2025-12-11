package main

import (
	"fmt"
	"log"

	"github.com/bit2swaz/codechef-recruit/backend/internal/game"
	"github.com/bit2swaz/codechef-recruit/backend/internal/store"
)

func main() {
	fmt.Println("=== Manual QA Test for Raja Mantri Chor Sipahi Game ===\n")

	// Test 1: Create a room with 4 hardcoded players
	fmt.Println("Test 1: Creating room with 4 players...")
	rm := store.NewRoomManager()
	room := rm.CreateRoom("QA-TEST")

	players := []struct {
		ID   string
		Name string
	}{
		{"player-1", "Alice"},
		{"player-2", "Bob"},
		{"player-3", "Charlie"},
		{"player-4", "David"},
	}

	for _, p := range players {
		player := store.Player{
			ID:    p.ID,
			Name:  p.Name,
			Role:  "", // No role initially
			Score: 0,
		}
		room.AddPlayer(player)
		fmt.Printf("  ✓ Added player: %s (ID: %s)\n", p.Name, p.ID)
	}

	fmt.Printf("\nRoom '%s' created with %d players\n", room.ID, len(room.GetPlayers()))
	fmt.Printf("Initial Status: %s\n\n", room.GetStatus())

	// Test 2: Assign roles randomly
	fmt.Println("Test 2: Assigning roles randomly...")
	game.AssignRoles(room)

	// Get updated players and display roles
	updatedPlayers := room.GetPlayers()
	fmt.Println("\n--- Role Assignments ---")

	roleCount := make(map[string]int)
	for _, player := range updatedPlayers {
		fmt.Printf("  %s (ID: %s) -> %s\n", player.Name, player.ID, player.Role)
		roleCount[player.Role]++
	}

	fmt.Println("\n--- Role Distribution ---")
	expectedRoles := []string{"Raja", "Mantri", "Chor", "Sipahi"}
	allRolesPresent := true

	for _, role := range expectedRoles {
		count := roleCount[role]
		status := "✓"
		if count != 1 {
			status = "✗"
			allRolesPresent = false
		}
		fmt.Printf("  %s %s: %d player(s)\n", status, role, count)
	}

	if allRolesPresent {
		fmt.Println("\n✓ SUCCESS: All roles are unique and present!")
	} else {
		fmt.Println("\n✗ FAILURE: Role distribution is incorrect!")
	}

	fmt.Printf("\nRoom Status after role assignment: %s\n", room.GetStatus())

	// Additional: Test ProcessGuess
	fmt.Println("\n--- Bonus: Testing ProcessGuess ---")

	// Find Mantri and Chor IDs
	var mantriID, chorID string
	for _, player := range updatedPlayers {
		if player.Role == "Mantri" {
			mantriID = player.ID
		} else if player.Role == "Chor" {
			chorID = player.ID
		}
	}

	// Test correct guess
	fmt.Println("\nScenario A: Mantri makes CORRECT guess...")
	result, err := game.ProcessGuess(room, mantriID, chorID)
	if err != nil {
		log.Fatalf("Error during guess: %v", err)
	}

	if result.Correct {
		fmt.Println("  ✓ Guess was CORRECT!")
	} else {
		fmt.Println("  ✗ Guess was WRONG!")
	}

	fmt.Println("\n--- Final Scores ---")
	finalPlayers := room.GetPlayers()
	totalScore := 0

	for _, player := range finalPlayers {
		fmt.Printf("  %s (%s): %d points\n", player.Name, player.Role, player.Score)
		totalScore += player.Score
	}

	fmt.Printf("\nTotal Score: %d (should be 2300)\n", totalScore)
	fmt.Printf("Final Room Status: %s\n", room.GetStatus())

	// Verify scoring rules
	fmt.Println("\n--- Scoring Verification ---")
	scoresByRole := make(map[string]int)
	for _, player := range finalPlayers {
		scoresByRole[player.Role] = player.Score
	}

	expectedScores := map[string]int{
		"Raja":   1000,
		"Mantri": 800, // Correct guess
		"Sipahi": 500,
		"Chor":   0,
	}

	allCorrect := true
	for role, expectedScore := range expectedScores {
		actualScore := scoresByRole[role]
		status := "✓"
		if actualScore != expectedScore {
			status = "✗"
			allCorrect = false
		}
		fmt.Printf("  %s %s: %d (expected %d)\n", status, role, actualScore, expectedScore)
	}

	if allCorrect {
		fmt.Println("\n✓ SUCCESS: All scores are correct!")
	} else {
		fmt.Println("\n✗ FAILURE: Some scores are incorrect!")
	}

	fmt.Println("\n=== Manual QA Test Complete ===")
}

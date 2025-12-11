package game

import (
	"testing"

	"github.com/bit2swaz/codechef-recruit/backend/internal/store"
)

// TestAssignRolesExactly4Players tests that roles are assigned when exactly 4 players are present
func TestAssignRolesExactly4Players(t *testing.T) {
	// Create a room manager and room
	rm := store.NewRoomManager()
	room := rm.CreateRoom("TEST1")

	// Add exactly 4 players
	players := []string{"Alice", "Bob", "Charlie", "David"}
	for _, name := range players {
		player := store.Player{
			ID:    name + "-id",
			Name:  name,
			Role:  "", // Initially no role
			Score: 0,
		}
		room.AddPlayer(player)
	}

	// Verify room has 4 players
	roomPlayers := room.GetPlayers()
	if len(roomPlayers) != 4 {
		t.Fatalf("Expected 4 players, got %d", len(roomPlayers))
	}

	// Call AssignRoles
	AssignRoles(room)

	// Verify roles were assigned
	updatedPlayers := room.GetPlayers()
	if len(updatedPlayers) != 4 {
		t.Fatalf("Expected 4 players after assignment, got %d", len(updatedPlayers))
	}

	// Verify each player has a role
	expectedRoles := map[string]bool{
		"Raja":   false,
		"Mantri": false,
		"Chor":   false,
		"Sipahi": false,
	}

	for _, player := range updatedPlayers {
		if player.Role == "" {
			t.Errorf("Player %s has no role assigned", player.Name)
		}

		// Check if role is valid
		if _, valid := expectedRoles[player.Role]; !valid {
			t.Errorf("Player %s has invalid role: %s", player.Name, player.Role)
		}

		// Mark role as used
		if expectedRoles[player.Role] {
			t.Errorf("Role %s assigned to multiple players", player.Role)
		}
		expectedRoles[player.Role] = true
	}

	// Verify all roles were assigned
	for role, assigned := range expectedRoles {
		if !assigned {
			t.Errorf("Role %s was not assigned to any player", role)
		}
	}

	// Verify room status changed to 'GUESSING'
	status := room.GetStatus()
	if status != "GUESSING" {
		t.Errorf("Expected room status to be 'GUESSING', got '%s'", status)
	}

	t.Log("Successfully assigned roles to all 4 players:")
	for _, player := range updatedPlayers {
		t.Logf("  %s -> %s", player.Name, player.Role)
	}
}

// TestAssignRolesLessThan4Players tests that roles are NOT assigned when less than 4 players
func TestAssignRolesLessThan4Players(t *testing.T) {
	rm := store.NewRoomManager()
	room := rm.CreateRoom("TEST2")

	// Add only 3 players
	for i := 1; i <= 3; i++ {
		player := store.Player{
			ID:    string(rune(i + 48)),
			Name:  "Player" + string(rune(i+48)),
			Role:  "",
			Score: 0,
		}
		room.AddPlayer(player)
	}

	// Call AssignRoles
	AssignRoles(room)

	// Verify roles were NOT assigned
	players := room.GetPlayers()
	for _, player := range players {
		if player.Role != "" {
			t.Errorf("Player %s should not have a role (only 3 players), but has: %s", player.Name, player.Role)
		}
	}

	// Verify status is still WAITING
	status := room.GetStatus()
	if status != "WAITING" {
		t.Errorf("Expected room status to remain 'WAITING', got '%s'", status)
	}
}

// TestAssignRolesMoreThan4Players tests that roles are NOT assigned when more than 4 players
func TestAssignRolesMoreThan4Players(t *testing.T) {
	rm := store.NewRoomManager()
	room := rm.CreateRoom("TEST3")

	// Add 5 players
	for i := 1; i <= 5; i++ {
		player := store.Player{
			ID:    string(rune(i + 48)),
			Name:  "Player" + string(rune(i+48)),
			Role:  "",
			Score: 0,
		}
		room.AddPlayer(player)
	}

	// Call AssignRoles
	AssignRoles(room)

	// Verify roles were NOT assigned
	players := room.GetPlayers()
	for _, player := range players {
		if player.Role != "" {
			t.Errorf("Player %s should not have a role (5 players), but has: %s", player.Name, player.Role)
		}
	}

	// Verify status is still WAITING
	status := room.GetStatus()
	if status != "WAITING" {
		t.Errorf("Expected room status to remain 'WAITING', got '%s'", status)
	}
}

// TestAssignRolesRandomness tests that roles are shuffled randomly
func TestAssignRolesRandomness(t *testing.T) {
	// Run multiple times to check for randomness
	roleAssignments := make(map[string]map[string]int) // player position -> role -> count

	iterations := 100
	for iter := 0; iter < iterations; iter++ {
		rm := store.NewRoomManager()
		room := rm.CreateRoom("TEST-RANDOM")

		// Add 4 players with consistent names
		players := []string{"Player0", "Player1", "Player2", "Player3"}
		for _, name := range players {
			player := store.Player{
				ID:    name + "-id",
				Name:  name,
				Role:  "",
				Score: 0,
			}
			room.AddPlayer(player)
		}

		// Assign roles
		AssignRoles(room)

		// Record role assignments by position
		updatedPlayers := room.GetPlayers()
		for i, player := range updatedPlayers {
			posKey := string(rune(i + 48))
			if roleAssignments[posKey] == nil {
				roleAssignments[posKey] = make(map[string]int)
			}
			roleAssignments[posKey][player.Role]++
		}
	}

	// Verify that each position got different roles across iterations
	// (indicates randomness)
	for pos, roles := range roleAssignments {
		if len(roles) < 2 {
			t.Errorf("Position %s only got %d different roles (expected variation for randomness)", pos, len(roles))
		}
		t.Logf("Position %s role distribution: %v", pos, roles)
	}
}

// TestAssignRolesWithEmptyRoom tests behavior with empty room
func TestAssignRolesWithEmptyRoom(t *testing.T) {
	rm := store.NewRoomManager()
	room := rm.CreateRoom("TEST-EMPTY")

	// Don't add any players

	// Call AssignRoles (should do nothing)
	AssignRoles(room)

	// Verify status is still WAITING
	status := room.GetStatus()
	if status != "WAITING" {
		t.Errorf("Expected room status to remain 'WAITING', got '%s'", status)
	}

	// Verify no players
	players := room.GetPlayers()
	if len(players) != 0 {
		t.Errorf("Expected 0 players, got %d", len(players))
	}
}

// TestAssignRolesThreadSafety tests concurrent role assignments
func TestAssignRolesThreadSafety(t *testing.T) {
	rm := store.NewRoomManager()

	// Create multiple rooms
	numRooms := 10
	rooms := make([]*store.Room, numRooms)

	for i := 0; i < numRooms; i++ {
		room := rm.CreateRoom("ROOM-" + string(rune(i+48)))

		// Add 4 players to each room
		for j := 0; j < 4; j++ {
			player := store.Player{
				ID:    "player-" + string(rune(i+48)) + "-" + string(rune(j+48)),
				Name:  "Player" + string(rune(j+48)),
				Role:  "",
				Score: 0,
			}
			room.AddPlayer(player)
		}

		rooms[i] = room
	}

	// Assign roles concurrently
	done := make(chan bool, numRooms)
	for _, room := range rooms {
		go func(r *store.Room) {
			AssignRoles(r)
			done <- true
		}(room)
	}

	// Wait for all to complete
	for i := 0; i < numRooms; i++ {
		<-done
	}

	// Verify all rooms have correct assignments
	for i, room := range rooms {
		players := room.GetPlayers()
		if len(players) != 4 {
			t.Errorf("Room %d: expected 4 players, got %d", i, len(players))
		}

		status := room.GetStatus()
		if status != "GUESSING" {
			t.Errorf("Room %d: expected status 'GUESSING', got '%s'", i, status)
		}

		// Verify unique roles
		roleCount := make(map[string]int)
		for _, player := range players {
			if player.Role == "" {
				t.Errorf("Room %d: Player %s has no role", i, player.Name)
			}
			roleCount[player.Role]++
		}

		if len(roleCount) != 4 {
			t.Errorf("Room %d: expected 4 unique roles, got %d", i, len(roleCount))
		}
	}
}

// TestProcessGuessCorrect tests when Mantri guesses the Chor correctly
func TestProcessGuessCorrect(t *testing.T) {
	rm := store.NewRoomManager()
	room := rm.CreateRoom("GUESS1")

	// Add 4 players with specific roles
	players := []store.Player{
		{ID: "raja-1", Name: "Raja Player", Role: "Raja", Score: 0},
		{ID: "mantri-1", Name: "Mantri Player", Role: "Mantri", Score: 0},
		{ID: "chor-1", Name: "Chor Player", Role: "Chor", Score: 0},
		{ID: "sipahi-1", Name: "Sipahi Player", Role: "Sipahi", Score: 0},
	}

	for _, player := range players {
		room.AddPlayer(player)
	}

	// Manually set roles and status to GUESSING
	room.UpdatePlayersAndStatus(players, "GUESSING")

	// Mantri makes correct guess
	result, err := ProcessGuess(room, "mantri-1", "chor-1")

	// Verify no error
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify guess was correct
	if !result.Correct {
		t.Error("Expected correct guess, got incorrect")
	}

	// Verify actual Chor ID
	if result.ActualChorID != "chor-1" {
		t.Errorf("Expected actual Chor ID to be chor-1, got %s", result.ActualChorID)
	}

	// Verify scores - Scenario A (Correct Guess)
	expectedScores := map[string]int{
		"raja-1":   1000, // Raja always gets 1000
		"mantri-1": 800,  // Mantri gets 800 on correct guess
		"sipahi-1": 500,  // Sipahi always gets 500
		"chor-1":   0,    // Chor gets 0 when caught
	}

	for playerID, expectedScore := range expectedScores {
		actualScore := result.UpdatedScores[playerID]
		if actualScore != expectedScore {
			t.Errorf("Player %s: expected score %d, got %d", playerID, expectedScore, actualScore)
		}
	}

	// Verify room status changed to FINISHED
	status := room.GetStatus()
	if status != "FINISHED" {
		t.Errorf("Expected room status to be FINISHED, got %s", status)
	}

	t.Logf("Correct guess - Final scores: Raja=%d, Mantri=%d, Sipahi=%d, Chor=%d",
		result.UpdatedScores["raja-1"],
		result.UpdatedScores["mantri-1"],
		result.UpdatedScores["sipahi-1"],
		result.UpdatedScores["chor-1"])
}

// TestProcessGuessWrong tests when Mantri guesses incorrectly
func TestProcessGuessWrong(t *testing.T) {
	rm := store.NewRoomManager()
	room := rm.CreateRoom("GUESS2")

	// Add 4 players with specific roles
	players := []store.Player{
		{ID: "raja-2", Name: "Raja Player", Role: "Raja", Score: 0},
		{ID: "mantri-2", Name: "Mantri Player", Role: "Mantri", Score: 0},
		{ID: "chor-2", Name: "Chor Player", Role: "Chor", Score: 0},
		{ID: "sipahi-2", Name: "Sipahi Player", Role: "Sipahi", Score: 0},
	}

	for _, player := range players {
		room.AddPlayer(player)
	}

	room.UpdatePlayersAndStatus(players, "GUESSING")

	// Mantri makes wrong guess (guesses Sipahi instead of Chor)
	result, err := ProcessGuess(room, "mantri-2", "sipahi-2")

	// Verify no error
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify guess was incorrect
	if result.Correct {
		t.Error("Expected incorrect guess, got correct")
	}

	// Verify actual Chor ID
	if result.ActualChorID != "chor-2" {
		t.Errorf("Expected actual Chor ID to be chor-2, got %s", result.ActualChorID)
	}

	// Verify scores - Scenario B (Wrong Guess)
	expectedScores := map[string]int{
		"raja-2":   1000, // Raja always gets 1000
		"mantri-2": 0,    // Mantri gets 0 on wrong guess
		"sipahi-2": 500,  // Sipahi always gets 500
		"chor-2":   800,  // Chor gets 800 (steals Mantri's points)
	}

	for playerID, expectedScore := range expectedScores {
		actualScore := result.UpdatedScores[playerID]
		if actualScore != expectedScore {
			t.Errorf("Player %s: expected score %d, got %d", playerID, expectedScore, actualScore)
		}
	}

	// Verify room status changed to FINISHED
	status := room.GetStatus()
	if status != "FINISHED" {
		t.Errorf("Expected room status to be FINISHED, got %s", status)
	}

	t.Logf("Wrong guess - Final scores: Raja=%d, Mantri=%d, Sipahi=%d, Chor=%d",
		result.UpdatedScores["raja-2"],
		result.UpdatedScores["mantri-2"],
		result.UpdatedScores["sipahi-2"],
		result.UpdatedScores["chor-2"])
}

// TestProcessGuessNonMantriCaller tests that non-Mantri cannot make guess
func TestProcessGuessNonMantriCaller(t *testing.T) {
	rm := store.NewRoomManager()
	room := rm.CreateRoom("GUESS3")

	players := []store.Player{
		{ID: "raja-3", Name: "Raja Player", Role: "Raja", Score: 0},
		{ID: "mantri-3", Name: "Mantri Player", Role: "Mantri", Score: 0},
		{ID: "chor-3", Name: "Chor Player", Role: "Chor", Score: 0},
		{ID: "sipahi-3", Name: "Sipahi Player", Role: "Sipahi", Score: 0},
	}

	for _, player := range players {
		room.AddPlayer(player)
	}

	room.UpdatePlayersAndStatus(players, "GUESSING")

	// Try to make guess as Raja (not Mantri)
	result, err := ProcessGuess(room, "raja-3", "chor-3")

	// Should return error
	if err == nil {
		t.Fatal("Expected error when non-Mantri makes guess, got nil")
	}

	if result != nil {
		t.Error("Expected nil result when error occurs")
	}

	// Verify error message
	expectedError := "only the Mantri can make a guess"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

	// Verify room status is still GUESSING (no changes)
	status := room.GetStatus()
	if status != "GUESSING" {
		t.Errorf("Expected room status to remain GUESSING, got %s", status)
	}
}

// TestProcessGuessNoRolesAssigned tests error when roles not assigned
func TestProcessGuessNoRolesAssigned(t *testing.T) {
	rm := store.NewRoomManager()
	room := rm.CreateRoom("GUESS4")

	// Add 4 players but without roles
	players := []store.Player{
		{ID: "player-1", Name: "Player 1", Role: "", Score: 0},
		{ID: "player-2", Name: "Player 2", Role: "", Score: 0},
		{ID: "player-3", Name: "Player 3", Role: "", Score: 0},
		{ID: "player-4", Name: "Player 4", Role: "", Score: 0},
	}

	for _, player := range players {
		room.AddPlayer(player)
	}

	// Try to process guess
	result, err := ProcessGuess(room, "player-1", "player-2")

	// Should return error
	if err == nil {
		t.Fatal("Expected error when roles not assigned, got nil")
	}

	if result != nil {
		t.Error("Expected nil result when error occurs")
	}

	// Verify error message
	expectedError := "room does not have all required roles assigned"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

// TestProcessGuessWithRealFlow tests the complete flow with AssignRoles
func TestProcessGuessWithRealFlow(t *testing.T) {
	rm := store.NewRoomManager()
	room := rm.CreateRoom("FLOW1")

	// Add 4 players
	playerNames := []string{"Alice", "Bob", "Charlie", "David"}
	for _, name := range playerNames {
		player := store.Player{
			ID:    name + "-id",
			Name:  name,
			Role:  "",
			Score: 0,
		}
		room.AddPlayer(player)
	}

	// Assign roles randomly
	AssignRoles(room)

	// Get players to find who is Mantri and Chor
	players := room.GetPlayers()
	var mantriID, chorID string

	for _, player := range players {
		if player.Role == "Mantri" {
			mantriID = player.ID
		} else if player.Role == "Chor" {
			chorID = player.ID
		}
	}

	// Mantri makes correct guess
	result, err := ProcessGuess(room, mantriID, chorID)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !result.Correct {
		t.Error("Expected correct guess")
	}

	// Verify total score is 2300 (1000 + 800 + 500 + 0)
	totalScore := 0
	for _, score := range result.UpdatedScores {
		totalScore += score
	}

	expectedTotal := 2300
	if totalScore != expectedTotal {
		t.Errorf("Expected total score to be %d, got %d", expectedTotal, totalScore)
	}

	// Verify status is FINISHED
	if room.GetStatus() != "FINISHED" {
		t.Errorf("Expected status FINISHED, got %s", room.GetStatus())
	}

	t.Logf("Complete flow test passed - Total score: %d", totalScore)
}

// TestProcessGuessScoreTotals tests score totals for both scenarios
func TestProcessGuessScoreTotals(t *testing.T) {
	scenarios := []struct {
		name          string
		correctGuess  bool
		expectedTotal int
	}{
		{"Correct Guess", true, 2300}, // 1000 + 800 + 500 + 0
		{"Wrong Guess", false, 2300},  // 1000 + 0 + 500 + 800
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			rm := store.NewRoomManager()
			room := rm.CreateRoom("TOTAL-" + scenario.name)

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

			// Make guess based on scenario
			var guessID string
			if scenario.correctGuess {
				guessID = "chor"
			} else {
				guessID = "sipahi" // Wrong guess
			}

			result, err := ProcessGuess(room, "mantri", guessID)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Calculate total
			total := 0
			for _, score := range result.UpdatedScores {
				total += score
			}

			if total != scenario.expectedTotal {
				t.Errorf("Expected total %d, got %d", scenario.expectedTotal, total)
			}
		})
	}
}

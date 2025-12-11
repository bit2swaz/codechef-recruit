package store

import (
	"fmt"
	"sync"
	"testing"
)

func TestCreateRoom(t *testing.T) {
	rm := NewRoomManager()

	roomID := "TEST"
	room := rm.CreateRoom(roomID)

	if room == nil {
		t.Fatal("CreateRoom returned nil")
	}

	if room.ID != roomID {
		t.Errorf("Expected room ID to be %s, got %s", roomID, room.ID)
	}

	if room.Status != "WAITING" {
		t.Errorf("Expected room status to be WAITING, got %s", room.Status)
	}

	if len(room.Players) != 0 {
		t.Errorf("Expected room to have 0 players, got %d", len(room.Players))
	}

	retrievedRoom := rm.GetRoom(roomID)
	if retrievedRoom == nil {
		t.Fatal("GetRoom returned nil for existing room")
	}

	if retrievedRoom.ID != roomID {
		t.Errorf("Expected retrieved room ID to be %s, got %s", roomID, retrievedRoom.ID)
	}

	nonExistentRoom := rm.GetRoom("NONEXISTENT")
	if nonExistentRoom != nil {
		t.Error("Expected GetRoom to return nil for non-existent room")
	}
}

func TestConcurrentJoins(t *testing.T) {
	rm := NewRoomManager()
	roomID := "CONC"
	room := rm.CreateRoom(roomID)

	numGoroutines := 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	errChan := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(playerNum int) {
			defer wg.Done()

			player := Player{
				ID:    fmt.Sprintf("player-%d", playerNum),
				Name:  fmt.Sprintf("Player %d", playerNum),
				Role:  "Player",
				Score: 0,
			}

			room.AddPlayer(player)

			players := room.GetPlayers()
			if len(players) == 0 {
				errChan <- fmt.Errorf("goroutine %d: expected at least 1 player", playerNum)
			}

			retrievedRoom := rm.GetRoom(roomID)
			if retrievedRoom == nil {
				errChan <- fmt.Errorf("goroutine %d: failed to retrieve room", playerNum)
			}
		}(i)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		t.Error(err)
	}

	finalPlayers := room.GetPlayers()
	if len(finalPlayers) != numGoroutines {
		t.Errorf("Expected %d players in room, got %d", numGoroutines, len(finalPlayers))
	}

	playerIDs := make(map[string]bool)
	for _, player := range finalPlayers {
		if playerIDs[player.ID] {
			t.Errorf("Duplicate player ID found: %s", player.ID)
		}
		playerIDs[player.ID] = true
	}

	if len(playerIDs) != numGoroutines {
		t.Errorf("Expected %d unique player IDs, got %d", numGoroutines, len(playerIDs))
	}
}

func TestConcurrentRoomOperations(t *testing.T) {
	rm := NewRoomManager()
	numOperations := 20

	var wg sync.WaitGroup
	wg.Add(numOperations)

	for i := 0; i < numOperations; i++ {
		go func(opNum int) {
			defer wg.Done()

			if opNum%2 == 0 {
				roomID := fmt.Sprintf("ROOM-%d", opNum)
				room := rm.CreateRoom(roomID)
				if room == nil {
					t.Errorf("Failed to create room %s", roomID)
				}
			} else {
				roomID := fmt.Sprintf("ROOM-%d", opNum-1)
				_ = rm.GetRoom(roomID)
			}
		}(i)
	}

	wg.Wait()

	expectedRooms := numOperations / 2
	actualRooms := 0

	for i := 0; i < numOperations; i += 2 {
		roomID := fmt.Sprintf("ROOM-%d", i)
		if rm.GetRoom(roomID) != nil {
			actualRooms++
		}
	}

	if actualRooms != expectedRooms {
		t.Errorf("Expected %d rooms, got %d", expectedRooms, actualRooms)
	}
}

func TestRoomStatusUpdates(t *testing.T) {
	rm := NewRoomManager()
	room := rm.CreateRoom("STATUS-TEST")

	numGoroutines := 5
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(routineNum int) {
			defer wg.Done()

			statuses := []string{"WAITING", "PLAYING", "FINISHED"}
			for _, status := range statuses {
				room.UpdateStatus(status)
				currentStatus := room.GetStatus()
				if currentStatus != "WAITING" && currentStatus != "PLAYING" && currentStatus != "FINISHED" {
					t.Errorf("Invalid status: %s", currentStatus)
				}
			}
		}(i)
	}

	wg.Wait()

	finalStatus := room.GetStatus()
	if finalStatus != "WAITING" && finalStatus != "PLAYING" && finalStatus != "FINISHED" {
		t.Errorf("Final status is invalid: %s", finalStatus)
	}
}

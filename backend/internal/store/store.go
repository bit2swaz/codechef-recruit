package store

import "sync"

type Player struct {
	ID    string
	Name  string
	Role  string
	Score int
}

type Room struct {
	ID      string
	Players []Player
	Status  string
	mu      sync.Mutex
}

type RoomManager struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*Room),
	}
}

func (rm *RoomManager) CreateRoom(id string) *Room {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	room := &Room{
		ID:      id,
		Players: make([]Player, 0),
		Status:  "WAITING",
	}

	rm.rooms[id] = room
	return room
}

func (rm *RoomManager) GetRoom(id string) *Room {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	return rm.rooms[id]
}

func (r *Room) AddPlayer(player Player) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Players = append(r.Players, player)
}

func (r *Room) UpdateStatus(status string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Status = status
}

func (r *Room) GetPlayers() []Player {
	r.mu.Lock()
	defer r.mu.Unlock()

	playersCopy := make([]Player, len(r.Players))
	copy(playersCopy, r.Players)
	return playersCopy
}

func (r *Room) GetStatus() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.Status
}

func (r *Room) UpdatePlayersAndStatus(players []Player, status string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Players = players
	r.Status = status
}

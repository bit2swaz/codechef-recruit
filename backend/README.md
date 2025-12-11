# CodeChef Recruit - Backend

A real-time multiplayer game backend built with Go, featuring WebSocket communication for live gameplay.

## Game Rules

**Raja-Mantri-Chor-Sipahi** is a classic Indian game for exactly 4 players.

### Roles
- **Raja (King)** - The ruler
- **Mantri (Minister)** - The advisor who must identify the Chor
- **Chor (Thief)** - Tries to hide among the players
- **Sipahi (Police)** - The enforcer

### Gameplay Flow

1. **Room Creation**: One player creates a room and becomes the admin
2. **Joining**: Three more players join (total 4 players required)
3. **Role Assignment**: Roles are randomly shuffled and assigned
4. **Guessing Phase**: The Mantri must identify who the Chor is
5. **Scoring**: Points are distributed based on whether the guess was correct

### Point Distribution

#### Scenario A: Mantri Guesses Correctly ✅
- **Raja**: 1000 points
- **Mantri**: 800 points
- **Sipahi**: 500 points
- **Chor**: 0 points

#### Scenario B: Mantri Guesses Incorrectly ❌
- **Raja**: 1000 points (always gets points)
- **Mantri**: 0 points (penalty for wrong guess)
- **Sipahi**: 500 points (always gets points)
- **Chor**: 800 points (steals Mantri's points)

## Technology Stack

- **Language**: Go (standard library)
- **HTTP Router**: [gorilla/mux](https://github.com/gorilla/mux) v1.8.1
- **WebSocket**: [gorilla/websocket](https://github.com/gorilla/websocket) v1.5.3
- **Concurrency**: Thread-safe with `sync.RWMutex`
- **Testing**: Table-driven tests with race detector

## Architecture

```
backend/
├── cmd/
│   ├── server/          # Main server entry point
│   ├── ws-client/       # WebSocket test client
│   └── qa-test/         # Manual QA utilities
├── internal/
│   ├── store/           # Thread-safe in-memory data store
│   ├── handlers/        # HTTP and WebSocket handlers
│   └── game/            # Game logic (roles, scoring)
└── go.mod
```

## API Endpoints

### Room Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/room/create` | Create a new room |
| POST | `/room/join` | Join an existing room |
| GET | `/room/{roomId}` | Get room details |

### Game Actions

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/game/start` | Start game (requires 4 players) |
| POST | `/game/guess` | Submit Mantri's guess |

### WebSocket

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/ws/{roomId}?playerId={playerId}` | Connect to room's WebSocket |

## API Examples

### 1. Create Room

```bash
curl -X POST http://localhost:8080/room/create \
  -H "Content-Type: application/json" \
  -d '{"playerName":"Alice"}'
```

**Response:**
```json
{
  "roomId": "ABCD"
}
```

### 2. Join Room

```bash
curl -X POST http://localhost:8080/room/join \
  -H "Content-Type: application/json" \
  -d '{
    "roomId": "ABCD",
    "playerName": "Bob"
  }'
```

**Response:**
```json
{
  "message": "Successfully joined room",
  "roomId": "ABCD"
}
```

### 3. Get Room Details

```bash
curl http://localhost:8080/room/ABCD
```

**Response:**
```json
{
  "roomId": "ABCD",
  "status": "WAITING",
  "players": [
    {
      "id": "20251211210336-ઐ",
      "name": "Alice",
      "score": 0
    },
    {
      "id": "20251211210336-Ὀ",
      "name": "Bob",
      "score": 0
    }
  ]
}
```

**Note**: Player roles are hidden until the game is finished.

### 4. Start Game

```bash
curl -X POST http://localhost:8080/game/start \
  -H "Content-Type: application/json" \
  -d '{"roomId":"ABCD"}'
```

**Response:**
```json
{
  "message": "Game started"
}
```

**Requirements**: Exactly 4 players must be in the room.

**WebSocket Broadcasts**:
- `GAME_START` - Sent to all players
- `YOUR_ROLE` - Sent privately to each player with their assigned role

### 5. Submit Guess

```bash
curl -X POST http://localhost:8080/game/guess \
  -H "Content-Type: application/json" \
  -d '{
    "roomId": "ABCD",
    "mantriPlayerId": "20251211210336-Ὀ",
    "guessedChorPlayerId": "20251211210336-ଧ"
  }'
```

**Response:**
```json
{
  "correct": true,
  "mantriId": "20251211210336-Ὀ",
  "chorId": "20251211210336-ଧ",
  "actualChorId": "20251211210336-ዂ",
  "updatedScores": {
    "20251211210336-ઐ": 1000,
    "20251211210336-Ὀ": 800,
    "20251211210336-ଧ": 500,
    "20251211210336-ዂ": 0
  }
}
```

**WebSocket Broadcasts**:
- `GUESS_RESULT` - Sent to all players with the outcome
- `GAME_END` - Sent to all players with final scores

### 6. WebSocket Connection

```javascript
// Connect to room's WebSocket
let ws = new WebSocket("ws://localhost:8080/ws/ABCD?playerId=20251211210336-ઐ");

ws.onopen = () => console.log("✅ Connected");
ws.onmessage = (event) => {
    let data = JSON.parse(event.data);
    console.log("📨 Received:", data);
};

// Example received messages:
// {"type":"connected","playerId":"...","roomId":"ABCD","data":{"message":"Connected to room"},"timestamp":1765467567}
// {"type":"PLAYER_JOINED","payload":{"name":"Bob","playerId":"..."}}
// {"type":"GAME_START","payload":{"message":"All players ready! Roles have been assigned."}}
// {"type":"YOUR_ROLE","payload":{"name":"Alice","role":"Raja"}}
// {"type":"GUESS_RESULT","payload":{"mantri":"Bob","correct":true,"scores":{...}}}
// {"type":"GAME_END","payload":{"message":"Game finished!","scores":{...}}}
```

## WebSocket Message Types

### Broadcast Messages (All Players)

**PLAYER_JOINED** - When a player joins the room
```json
{
  "type": "PLAYER_JOINED",
  "payload": {
    "name": "Bob",
    "playerId": "20251211210336-Ὀ"
  }
}
```

**GAME_START** - When roles are assigned
```json
{
  "type": "GAME_START",
  "payload": {
    "message": "All players ready! Roles have been assigned."
  }
}
```

**GUESS_RESULT** - When Mantri makes a guess
```json
{
  "type": "GUESS_RESULT",
  "payload": {
    "mantri": "Bob",
    "correct": true,
    "scores": {
      "player1": 1000,
      "player2": 800,
      "player3": 500,
      "player4": 0
    }
  }
}
```

**GAME_END** - When game finishes
```json
{
  "type": "GAME_END",
  "payload": {
    "message": "Game finished!",
    "scores": {
      "Alice": 1000,
      "Bob": 800,
      "Charlie": 500,
      "Diana": 0
    }
  }
}
```

### Private Messages (Single Player)

**YOUR_ROLE** - Sent privately to each player
```json
{
  "type": "YOUR_ROLE",
  "payload": {
    "name": "Alice",
    "role": "Raja"
  }
}
```

## How to Run

### Prerequisites

- Go 1.21 or higher
- Git

### Installation

```bash
# Clone the repository
git clone https://github.com/bit2swaz/codechef-recruit.git
cd codechef-recruit/backend

# Download dependencies
go mod download
```

### Run the Server

```bash
# Method 1: Run directly
go run cmd/server/main.go

# Method 2: Build and run
go build -o server cmd/server/main.go
./server
```

The server will start on `http://localhost:8080`

**Console Output:**
```
2025/12/11 21:02:28 WebSocket Hub initialized and running
2025/12/11 21:02:28 Server starting on port :8080
2025/12/11 21:02:28 WebSocket endpoint: ws://localhost:8080/ws/{{roomId}}?playerId={{playerId}}
```

### Run Tests

```bash
# Run all tests
go test -v ./...

# Run tests with race detector
go test -race -v ./...

# Run specific test package
go test -v ./internal/store/
go test -v ./internal/game/
go test -v ./cmd/server/

# Run WebSocket integration test
./test-broadcast-integration.sh
```

### Run Test WebSocket Client

```bash
# Build the client
go build -o ws-client cmd/ws-client/main.go

# Connect to a room
./ws-client -room=ABCD -player=alice-123

# Open multiple terminals for multiple clients
./ws-client -room=ABCD -player=bob-456
./ws-client -room=ABCD -player=charlie-789
```

## Development

### Project Structure

- **`cmd/server/main.go`** - Server entry point, route configuration
- **`internal/store/`** - Thread-safe in-memory data structures
  - `Room`, `Player`, `RoomManager`
- **`internal/handlers/`** - HTTP and WebSocket handlers
  - `room.go` - Room creation/joining/retrieval
  - `game.go` - Game start and guess submission
  - `websocket.go` - WebSocket connection handler
  - `websocket_hub.go` - WebSocket hub for managing connections
  - `broadcast.go` - Broadcast helper functions
- **`internal/game/`** - Game logic
  - `roles.go` - Role assignment and guess processing

### Adding New Features

1. **New API Endpoint**: Add handler in `internal/handlers/`
2. **New Game Logic**: Add function in `internal/game/`
3. **New Broadcast Event**: Add helper in `internal/handlers/broadcast.go`
4. **Register Route**: Update `cmd/server/main.go`

### Testing Guidelines

- Write table-driven tests for game logic
- Use `httptest.NewRecorder` for HTTP handler tests
- Use `httptest.NewServer` for WebSocket tests
- Always run with `-race` flag to detect race conditions
- Mock WebSocket clients for integration tests

## Configuration

### Port Configuration

Edit `cmd/server/main.go`:

```go
port := ":8080"  // Change to your desired port
```

### CORS Configuration

WebSocket upgrader allows all origins for development. For production, update `internal/handlers/websocket.go`:

```go
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        // Add your origin validation logic here
        return true  // Currently allows all origins
    },
}
```

### Ping/Pong Timeouts

Edit `internal/handlers/websocket.go`:

```go
const (
    writeWait = 10 * time.Second
    pongWait = 60 * time.Second      // Change timeout here
    pingPeriod = (pongWait * 9) / 10
    maxMessageSize = 512
)
```

## Thread Safety

The backend is fully thread-safe:

- **RoomManager** uses `sync.RWMutex` for room access
- **Room** uses `sync.Mutex` for player updates
- **Hub** uses `sync.RWMutex` for WebSocket client management
- All tests pass with Go's race detector

## Performance

- **In-Memory Storage**: Fast read/write operations
- **Goroutines**: Each WebSocket connection runs in separate goroutines
- **Non-Blocking Broadcasts**: Uses buffered channels to prevent blocking
- **Connection Limits**: No artificial limits, scales with system resources

## Error Handling

### HTTP Errors

- `400 Bad Request` - Invalid JSON, missing fields, or invalid game state
- `404 Not Found` - Room doesn't exist
- `500 Internal Server Error` - Server-side errors (logged)

### WebSocket Errors

- Connection rejected if:
  - `roomId` is missing
  - `playerId` query parameter is missing
  - Room doesn't exist

## Logging

Server logs all important events:

```
2025/12/11 21:03:36 Broadcast to room ABCD: type=PLAYER_JOINED, payload=map[name:Alice playerId:...]
2025/12/11 21:03:36 Client registered to room ABCD (PlayerID: alice-id). Total clients in room: 1
2025/12/11 21:03:38 Sent to player alice-id in room ABCD: type=YOUR_ROLE
2025/12/11 21:03:40 Client unregistered from room ABCD (PlayerID: alice-id). Remaining clients: 0
2025/12/11 21:03:40 Room ABCD is now empty and removed
```

## Documentation

- **`BROADCAST_INTEGRATION.md`** - Detailed broadcast implementation guide
- **`WEBSOCKET_TESTING.md`** - Manual testing guide with browser console examples
- **`test-broadcast-integration.sh`** - Integration test script

## Troubleshooting

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

### WebSocket Connection Failed

- Ensure room exists before connecting
- Include `playerId` query parameter
- Check server logs for detailed error messages
- Verify WebSocket URL format: `ws://localhost:8080/ws/{roomId}?playerId={playerId}`

### Broadcast Not Received

- Ensure WebSocket connection is established (check `ws.readyState === 1`)
- Verify you're connected to the correct room
- Check server logs to confirm broadcast was sent
- Try reconnecting the WebSocket

### Tests Failing

```bash
# Clean build cache
go clean -testcache

# Run tests individually
go test -v -run TestName ./...

# Check for race conditions
go test -race -v ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Write tests for your changes
4. Ensure all tests pass with race detector
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## License

This project is part of CodeChef recruitment assessment.

## Support

For issues or questions:
- Create an issue on GitHub
- Check existing documentation in `/backend/*.md` files
- Review test files for usage examples

---

Built with ❤️ using Go and WebSockets

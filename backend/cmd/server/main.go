package main

import (
	"log"
	"net/http"

	"github.com/bit2swaz/codechef-recruit/backend/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	handlers.InitHub()

	r := mux.NewRouter()

	r.HandleFunc("/room/create", handlers.CreateRoom).Methods("POST")
	r.HandleFunc("/room/join", handlers.JoinRoom).Methods("POST")
	r.HandleFunc("/room/{roomId}", handlers.GetRoom).Methods("GET")
	r.HandleFunc("/game/start", handlers.StartGame).Methods("POST")
	r.HandleFunc("/game/guess", handlers.SubmitGuess).Methods("POST")

	r.HandleFunc("/ws/{roomId}", handlers.HandleWebSocket).Methods("GET")

	port := ":8080"
	log.Printf("Server starting on port %s", port)
	log.Printf("WebSocket endpoint: ws://localhost%s/ws/{{roomId}}?playerId={{playerId}}", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}

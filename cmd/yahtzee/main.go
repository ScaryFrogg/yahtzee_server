package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ScaryFrogg/yahtzee_server/internal/service"
	"github.com/ScaryFrogg/yahtzee_server/internal/types"
	"github.com/gorilla/websocket"
)

type webSocketHandler struct {
	upgrader websocket.Upgrader
	room     *types.Room
}

func (wsh webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	playerID := r.URL.Query().Get("playerId")
	if playerID == "" {
		http.Error(w, "player parameter required", http.StatusBadRequest)
		return
	}

	conn, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error %s when upgrading connection to websocket", err)
		return
	}
	defer conn.Close()

	player, err := wsh.room.AddPlayer(playerID, conn)
	if err != nil {
		conn.WriteJSON(map[string]string{"error": err.Error()})
		return
	}
	board := player.Board
	for {
		var req types.Message
		if err := conn.ReadJSON(&req); err != nil {
			log.Printf("Failed to decode request: %v", err)
			return
		}
		log.Printf("Received message: %s\n", req)
		//TODO cleanup, refactor, split responsibility
		switch req.Type {
		case types.TypeSync:
			conn.WriteJSON(board)
		case types.TypeRoll:
			service.Roll(board, [6]bool{})
		case types.TypeReRoll:
			var payload types.ReRollPayload
			if err := json.Unmarshal(req.Payload, &payload); err != nil {
				log.Println(err)
				conn.WriteJSON("{'status':'failed_unmarshalling'}")
				return
			}
			service.Roll(board, payload.Changes)
			conn.WriteJSON(board.CurrentRoll)
			conn.WriteJSON(service.Calculate(board))
		case types.TypeCommit:
			var payload types.CommitPayload
			if err := json.Unmarshal(req.Payload, &payload); err != nil {
				log.Println(err)
				conn.WriteJSON("{'status':'failed_unmarshalling'}")
				return
			}
			service.Commit(board, payload.CommitIndex)
			conn.WriteJSON(board.Rows)
		default:
			conn.WriteJSON("{'status':'unknown_type'}")
		}
		for _, player := range wsh.room.Players {
			types.LogPlayerBoard(player)
		}
	}
}

func main() {
	webSocketHandler := webSocketHandler{
		room: types.CreateRoom("asd"),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
	http.Handle("/", webSocketHandler)
	log.Print("Starting server...")
	log.Fatal(http.ListenAndServe("localhost:3005", nil))
}

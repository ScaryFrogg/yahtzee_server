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
		//TODO cleanup, refactor, split responsibility
		switch req.Type {
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
		case types.TypeCommit:
			var payload types.CommitPayload
			if err := json.Unmarshal(req.Payload, &payload); err != nil {
				log.Println(err)
				conn.WriteJSON("{'status':'failed_unmarshalling'}")
				return
			}
			service.Commit(board, payload.CommitIndex)

			//start next round if everyone is ready
			if wsh.room.CheckAllCommitted() {
				//fan out to notify all players
				for _, p := range wsh.room.Players {
					service.Roll(p.Board, [6]bool{})
					p.Board.Waiting = true
					p.Conn.WriteJSON(p.Board.CurrentRoll)
				}
			}

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

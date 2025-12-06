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
}

func (wsh webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error %s when upgrading connection to websocket", err)
		return
	}
	defer conn.Close()
	//TODO handle close/disconnect
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
			conn.WriteJSON(testBoard)
		case types.TypeRoll:
			service.Roll(testBoard, [6]bool{})
		case types.TypeReRoll:
			var payload types.ReRollPayload
			if err := json.Unmarshal(req.Payload, &payload); err != nil {
				log.Println(err)
				conn.WriteJSON("{'status':'failed_unmarshalling'}")
				return
			}
			service.Roll(testBoard, payload.Changes)
			conn.WriteJSON(testBoard.CurrentRoll)
			conn.WriteJSON(service.Calculate(*testBoard))
		case types.TypeCommit:
			var payload types.CommitPayload
			if err := json.Unmarshal(req.Payload, &payload); err != nil {
				log.Println(err)
				conn.WriteJSON("{'status':'failed_unmarshalling'}")
				return
			}
			service.Commit(testBoard, payload.CommitIndex)
		default:
			conn.WriteJSON("{'status':'unknown_type'}")
		}

	}
}

var testBoard = types.NewBoard()

func main() {
	webSocketHandler := webSocketHandler{
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

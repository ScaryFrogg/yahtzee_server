package main

import (
	"log"
	"math/rand/v2"
	"net/http"

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
		case types.TypeRoll:
			arr := []int{rand.IntN(6) + 1, rand.IntN(6) + 1, rand.IntN(6) + 1, rand.IntN(6) + 1, rand.IntN(6) + 1, rand.IntN(6) + 1}
			conn.WriteJSON(arr)
		case types.TypeReRoll:
			// var p ReRollPayload
			// if err := json.Unmarshal(req.Payload, &p); err != nil {
			// 	return err
			// }
			conn.WriteJSON("{'status':'ok'}")
		default:
			conn.WriteJSON("{'status':'unknown_type'}")
		}
	}
}

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

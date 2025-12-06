package types

import "github.com/gorilla/websocket"

type Player struct {
	Id    string
	Conn  *websocket.Conn
	Board *Board
}

type Room struct {
	Id      string
	Players map[string]*Player
}

// TODO CLEANUP move to service
func CreateRoom(roomId string) *Room {
	return &Room{
		Id:      roomId,
		Players: make(map[string]*Player),
	}
}

func (room *Room) AddPlayer(playerID string, conn *websocket.Conn) (*Player, error) {
	player := &Player{
		Id:    playerID,
		Conn:  conn,
		Board: NewBoard(),
	}
	room.Players[playerID] = player
	return player, nil
}

func (room *Room) CheckAllCommitted() bool {
	for _, player := range room.Players {
		if !player.Board.Waiting {
			return false
		}
	}
	return len(room.Players) > 0
}

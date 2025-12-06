package types

import "encoding/json"

type Message struct {
	PlayerId string          `json:"playerId"`
	Type     MessageType     `json:"type"`
	Payload  json.RawMessage `json:"payload"`
}

type MessageType string

const (
	TypeRoll   MessageType = "roll"
	TypeReRoll MessageType = "reroll"
	TypeSync   MessageType = "sync"
	TypeCommit MessageType = "commit"
)

type ReRollPayload struct {
	Changes [6]bool `json:"changes"`
}

type CommitPayload struct {
	CommitIndex int `json:"commitIndex"`
}

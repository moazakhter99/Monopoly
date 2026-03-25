package models

import "encoding/json"



type WSMessage struct {
	Type    string          `json:"type"`    // e.g., "ROLL_DICE", "BUY_PROPERTY", "CHAT"
	Payload json.RawMessage `json:"payload"` // The specific data for that action
}

type WsError struct {
	Message string
	WsError int
}
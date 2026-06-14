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

type Request struct {
	PlayerId string `json:"player_id"`
	GameId string `json:"game_id"`

}

type ReqDiceRoll struct {
	PlayerId string `json:"player_id"`
	GameId string `json:"game_id"`
}

type RespDiceRoll struct {
	PlayerId string `json:"player_id"`
	GameId string `json:"game_id"`
	DiceVal int `json:"dice_val"`
}

type ReqMovePos struct {
	PlayerId string `json:"player_id"`
	GameId string `json:"game_id"`
	UpdateBy int `json:"update_by"`
	
}

type RespMovePos struct {
	PlayerId string `json:"player_id"`
	GameId string `json:"game_id"`
	NewPos int `json:"new_pos"`
}

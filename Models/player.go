package models


type Player struct {
	PlayerId 	string `json:"playerId"`
	Name 		string `json:"playerName"`
	Pos 		string `json:"pos"`
	Cash 		string `json:"cash"`
	Status 		string `json:"status"`
}
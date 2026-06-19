package models


type Player struct {
	PlayerId 	string `json:"playerId"`
	GameId string `json:"gameId"`
	Name 		string `json:"playerName"`
	Pos 		int `json:"pos"`
	Cash 		int `json:"cash"`
	Seq 		int `json:"seq"`
	Status 		string `json:"status"`
}
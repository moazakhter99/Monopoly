package models



type BlockInfo struct {
	InfoId	string	`json:"info_id"`
	Info	string	`json:"info"`
}

type Client struct {
	PlayerId string `json:"player_id"`
	GameId   string `json:"game_id"`
}

package models



type ReqCreateGame struct {
	MsgId 		string 	`json:"msgId"`
	GameId 		string 	`json:"gameId"`
	MatchId 	string	`json:"matchId"`
	Player 		*Player `json:"player"`
	Timestamp 	string 	`json:"timestamp"`
	// Game 	*Game 	`json:"game"`

}

type RespCreateGame struct {
	MsgId 		string 		 `json:"msgId"`
	GameId 		string 		 `json:"gameId"`
	Message 	string 		 `json:"message"`
	Timestamp 	string 		 `json:"timestamp"`
	GeneralResp *GeneralResp `json:"resp"`

}
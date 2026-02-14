package models


type ReqJoinGame struct {
	MsgId 		string `json:"msgId"`
	MatchId 	string `json:"matchId"`
	Player 		*Player `json:"player"`
	Timestamp 	string `json:"timestamp"`

}


type RespJoinGame struct {
	MsgId 		string `json:"msgId"`
	GameId 		string `json:"gameId"`
	Joined		string `json:"joined"`
	GeneralResp *GeneralResp `json:"resp"`

}

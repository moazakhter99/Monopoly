package models


type ReqStartGame struct {
	MsgId string
	GameId string

}



type RespStartGame struct {
	MsgId string
	GameId string
	Game *Game

	GeneralResp *GeneralResp `json:"resp"`
	
}
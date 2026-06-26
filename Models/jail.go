package models



type ReqJailInfo struct {
	MsgId		string		`json:"msg_id"`
	
}

type RespJailInfo struct {
	MsgId		string			`json:"msg_id"`
	BlockType	string			`json:"block_type"`
	BlockInfo	[]*BlockInfo	`json:"block_info"`
	
}
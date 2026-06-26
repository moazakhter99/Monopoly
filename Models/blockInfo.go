package models



type ReqBlockInfo struct {
	MsgId		string		`json:"msg_id"`
	
}

type RespBlockInfo struct {
	MsgId		string			`json:"msg_id"`
	BlockType	string			`json:"block_type"`
	BlockInfo	[]*BlockInfo	`json:"block_info"`
	
}
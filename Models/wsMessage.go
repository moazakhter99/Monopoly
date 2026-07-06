package models

import "encoding/json"

type WSMessage struct {
	Type    string          `json:"type"`    // e.g., "ROLL_DICE", "BUY_PROPERTY", "CHAT"
	Client	*Client			`json:"client"`
	Payload json.RawMessage `json:"payload"` // The specific data for that action
}

type WsError struct {
	Message string
	WsError int
}

type Request struct {
}

type RespDiceRoll struct {
	DiceVal  int    `json:"dice_val"`
}

type ReqMovePos struct {
	UpdateBy int    `json:"update_by"`
}

type RespMovePos struct {
	NewPos	int		`json:"new_pos"`
	BlockId	string	`json:"block_id"`
	State	string	`json:"sold"`
	Type	string	`json:"type"`
	OwnerId	string	`json:"owner_id"`
}

type ReqBuyBlock struct {
	BlockId  string `json:"block_id"`
	Price    int    `json:"price"`
	Colour   string `json:"colour"`
	CardNo   string `json:"card_no"`
	Seq		 int	`json:"seq"`
	Count	 int	`json:"count"`
}

type RespBuyBlock struct {
	BlockId  string `json:"block_id"`
	Buy      bool   `json:"buy"`
	Cash     int    `json:"cash"`
	Seq		 int	`json:"seq"`
	ChangePlayer bool `json:"change_player"`
}

type ReqActionCard struct {
	BlockId  string `json:"block_id"`
	CardId   string `json:"card_id"`
	Price    int    `json:"price"`
	Cash     int    `json:"cash"`
	Pos      int    `json:"position"`
	Type     string `json:"type"`
}

type RespActionCard struct {
	Cash     int    `json:"cash"`
	Pos      int    `json:"position"`
	InJail	 bool	`json:"in_jail"`
}

type Jail struct {
	InfoId string `json:"info_id"`
	Info   string `json:"info"`
}

type RespJail struct {
	GetOutCard bool   `json:"get_out_card"`
	NewPos     int    `json:"new_pos"`
	Jail       []Jail `json:"jail"`
	InJail     bool   `json:"In_jail"`
	Cash	   int	  `json:"cash"`
}

type ReqJail struct {
	InJail		bool	`json:"In_jail"`
	JailId		string	`json:"jail_id"`
	Cash		int		`json:"cash"`
	
}

type RespChangePlayer struct {
	NextPlayer	string	`json:"next_player"`
	Playing		bool	`json:"playing"`

}

type ReqCalculateRent struct {
	BlockId		string	`json:"block_id"`
	BlockType	string	`json:"block_type"`
	Price		int		`json:"price"`
	OwnerId		string	`json:"owner_id"`

}

type RespCalculateRent struct {
	BlockId		string	`json:"block_id"`
	BlockType	string	`json:"block_type"`
	OwnerId		string	`json:"owner_id"`
	RenterId	string	`json:"renter_id"`
	Rent		int		`json:"rent"`
}

type ReqPayRent struct {
	Cash		int		`json:"cash"`
	Rent		int		`json:"rent"`
	OwnerId		string	`json:"owner_id"`
}

type RespPayRent struct {
	Cash		int		`json:"cash"`
	Paid		bool	`json:"paid"`
	OwnerId		string	`json:"owner_id"`
	RenterId	string	`json:"renter_id"`
}

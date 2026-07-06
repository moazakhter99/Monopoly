package service

import (
	db "Monopoly/DB"
	gameplay "Monopoly/Gameplay"
	models "Monopoly/Models"
	"Monopoly/logger"

	"go.uber.org/zap"
)



type GameHub struct {
	Id string
	logger *zap.SugaredLogger
	// clinet *Client
	ReadMsg chan models.WSMessage
	WriteMsg chan models.WSMessage
	ClientMap map[string]*Client
	Register chan *Client
	Deregister chan *Client
	db *db.DbOperations

}


func CreateNewGameHub(logger *zap.SugaredLogger, db db.DbOperations) *GameHub {
	return &GameHub{
		Id: "hub123",
		logger: logger,
		ReadMsg: make(chan models.WSMessage, 2),
		WriteMsg: make(chan models.WSMessage, 2),
		ClientMap: map[string]*Client{},
		Register: make(chan *Client),
		Deregister: make(chan *Client),
		db: &db,
	}
}

func (h *GameHub) ProcessEvent(message any) {
	logger.ZapLogger.Infoln("Start ProcessEvent")

	var respMsg models.WSMessage
	var respByte []byte
	var respByteMap map[string][]byte
	all := true
	specific := false
	diff := false

	wsMsg := message.(models.WSMessage)
	h.logger.Infow("ReadMsg", "Type", wsMsg.Type, "Message", string(wsMsg.Payload))

	switch wsMsg.Type {

	case models.ROLLDICE:

		respByte = gameplay.RollDice(wsMsg.Payload, wsMsg.Client, h.ReadMsg)

	case models.MOVEPOS:

		respByte = gameplay.MovePos(wsMsg.Payload, wsMsg.Client, *h.db)

	case models.BUYBLOCK:

		respByte = gameplay.BuyBlock(wsMsg.Payload, wsMsg.Client, *h.db, h.ReadMsg)

	case models.CALCULATERENT:

		respByteMap = gameplay.CalculateRent(wsMsg.Payload, wsMsg.Client, *h.db)
		diff = true
		all = false

	case models.PAYRENT:

		respByteMap = gameplay.PayRent(wsMsg.Payload, wsMsg.Client, *h.db)
		diff = true
		all = false

	case models.SELL:
	
	case models.ACTIONCARD:
		
		respByte = gameplay.ActionCard(wsMsg.Payload, wsMsg.Client, *h.db, h.ReadMsg)

	case models.JAIL:

		respByte = gameplay.Jail(wsMsg.Payload, wsMsg.Client, *h.db)

	case models.CHANGEPLAYER:

		respByteMap = gameplay.ChangePlayer(wsMsg.Payload, wsMsg.Client, *h.db)
		diff = true
		all = false
	
	}

	respMsg.Type = wsMsg.Type

	if all {
		respMsg.Payload = respByte
		for _, cl := range h.ClientMap {
			cl.WriteMsg <- respMsg

		}
	}

	if specific {

	}

	if diff {
		for id, respMsgByte := range respByteMap {
			respMsg.Payload = respMsgByte
			cl := h.ClientMap[id]
			cl.WriteMsg <- respMsg
		}

	}

	logger.ZapLogger.Infoln("Exit ProcessEvent")
}


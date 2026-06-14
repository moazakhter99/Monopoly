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
	all := true

	wsMsg := message.(models.WSMessage)
	h.logger.Infow("ReadMsg", "Type", wsMsg.Type, "Message", string(wsMsg.Payload))

	switch wsMsg.Type {

	case models.ROLLDICE:

		respByte = gameplay.RollDice(wsMsg.Payload, h.ReadMsg)

	case models.MOVEPOS:

		respByte = gameplay.MovePos(wsMsg.Payload, *h.db)

	case models.BUY:
		return

	case models.SELL:
		return
	
	}

	respMsg.Payload = respByte
	respMsg.Type = wsMsg.Type

	if all {
		for _, cl := range h.ClientMap {
			cl.WriteMsg <- respMsg

		}
	}


	logger.ZapLogger.Infoln("Exit ProcessEvent")
}


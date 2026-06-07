package service

import (
	// models "Monopoly/Models"
	// client "Monopoly/Client"
	gameplay "Monopoly/Gameplay"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"

	"go.uber.org/zap"
)



type GameHub struct {
	logger *zap.SugaredLogger
	// clinet *Client
	ReadMsg chan models.WSMessage
	WriteMsg chan models.WSMessage
	ClientMap map[string]*Client
	Register chan *Client
	Deregister chan *Client

}


func CreateNewGameHub(logger *zap.SugaredLogger) *GameHub {
	return &GameHub{
		logger: logger,
		ReadMsg: make(chan models.WSMessage),
		WriteMsg: make(chan models.WSMessage),
		ClientMap: map[string]*Client{},
		Register: make(chan *Client),
		Deregister: make(chan *Client),
	}
}

func (h *GameHub) ProcessEvent(message any) {
	logger.ZapLogger.Infoln("Start ProcessEvent")

	var respMsg models.WSMessage
	var respByte []byte
	var playerId string

	wsMsg := message.(models.WSMessage)
	h.logger.Infow("ReadMsg", "Type", wsMsg.Type, "Message", string(wsMsg.Payload))

	switch wsMsg.Type {

	case models.ROLLDICE:
		logger.ZapLogger.Infoln("Roll Dice")
		var req models.ReqDiceRoll

		err := json.Unmarshal(wsMsg.Payload, &req)
		if err != nil {
			logger.ZapLogger.Errorw("UnMarshal Error", "Errror", err)
		}
		playerId = req.PlayerId

		diceVal := gameplay.DiceRoll()

		resp := models.RespDiceRoll{
			PlayerId: playerId,
			GameId: req.GameId,
			DiceVal: diceVal,
		}
		respByte, _ = json.Marshal(resp)

	case models.MOVEPOS:
		return

	case models.BUY:
		return

	case models.SELL:
		return
	
	}

	respMsg.Payload = respByte
	respMsg.Type = wsMsg.Type

	client := h.ClientMap[playerId]
	client.WriteMsg <- respMsg

	logger.ZapLogger.Infoln("Exit ProcessEvent")
}


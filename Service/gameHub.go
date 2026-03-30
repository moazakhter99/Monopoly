package service

import (
	// models "Monopoly/Models"
	// client "Monopoly/Client"
	models "Monopoly/Models"
	"Monopoly/logger"

	"go.uber.org/zap"
)



type GameHub struct {
	logger *zap.SugaredLogger
	clinet *Client
	ReadMsg chan models.WSMessage
	WriteMsg chan models.WSMessage

}


func CreateNewGameHub(logger *zap.SugaredLogger) *GameHub {
	return &GameHub{
		logger: logger,
		ReadMsg: make(chan models.WSMessage),
		WriteMsg: make(chan models.WSMessage),
	}
}

func (h *GameHub) ProcessEvent(message any) {
	logger.ZapLogger.Infoln("Start ProcessEvent")

	var respMsg models.WSMessage
	wsMsg := message.(models.WSMessage)
	
	h.logger.Infoln("ReadMsg", "Message", wsMsg)
	

	h.WriteMsg <- respMsg

	logger.ZapLogger.Infoln("Exit ProcessEvent")
}


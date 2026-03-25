package service

import (
	// models "Monopoly/Models"
	// client "Monopoly/Client"
	"Monopoly/logger"

	"go.uber.org/zap"
)



type GameHub struct {
	logger *zap.SugaredLogger
	clinet *Client

}


func CreateNewGameHub(logger *zap.SugaredLogger) *GameHub {
	return &GameHub{
		logger: logger,
	}
}

func (h *GameHub) ProcessEvent(message any) {
	logger.ZapLogger.Infoln("Start ProcessEvent")

	// wsMsg := message.(models.WSMessage)
	
	msg, _ := <- h.clinet.ReadMsg
	h.logger.Infoln("ReadMsg", "Event", msg.Type)
	


	// logger.ZapLogger.Infoln("Exit ProcessEvent")
}


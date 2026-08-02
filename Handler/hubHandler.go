package handler

import (
	db "Monopoly/DB"
	gameplay "Monopoly/Gameplay"
	models "Monopoly/Models"
	"Monopoly/logger"
	"Monopoly/router"
	"encoding/json"

	"go.uber.org/zap"
)


type NewHub struct {
	Id string
	logger *zap.SugaredLogger
	ReadMsg chan models.WSMessage
	WriteMsg chan models.WSMessage
	// ClientMap map[string]*service.Client
	// Register chan *service.Client
	// Deregister chan *service.Client
	db db.DbOperations
}

type HubController interface {
	HandleHub(msg []byte, wrChan chan []byte)

}

type GameHubController struct {
	game gameplay.Game
}

func CreateHubContoller(g gameplay.Game) *GameHubController {
	return &GameHubController{
		game: g,
	}
}

func (hc *GameHubController) HandleHub(req router.Request, wrChan chan []byte) {

	// wsMsg := message.(models.WSMessage)
	// hub.Infow("ReadMsg", "Type", wsMsg.Type, "Message", string(wsMsg.Payload))

	action := req.Action
	msg := req.Msg

	payload, err := hc.game.Validate(msg)
	if err != nil {
		logger.ZapLogger.Errorw("Validation Error", "Error", err)
		errResponse := models.RespError{
			Stage: models.VALIDATION,
			Error: err.Error(),
			Status: models.FAILED,
		}
		errResp, err := json.Marshal(errResponse)
		if err != nil {
			logger.ZapLogger.Errorw("JSON Error", "Error", err)
			return
		}
		wsMessage := models.WSMessage{
			Type: action,
			Payload: errResp,
			
		}
		resp, err := json.Marshal(wsMessage)
		if err != nil {
			logger.ZapLogger.Errorw("JSON Error", "Error", err)
			return
		}
		wrChan <- resp
		return
	}

	resp, err := hc.game.Play(payload)
	if err != nil {
		logger.ZapLogger.Errorw("Play Error", "Error", err)
		errResponse := models.RespError{
			Stage: models.PLAY,
			Error: err.Error(),
			Status: models.FAILED,
		}
		errResp, err := json.Marshal(errResponse)
		if err != nil {
			logger.ZapLogger.Errorw("JSON Error", "Error", err)
			return
		}
		wsMessage := models.WSMessage{
			Type: action,
			Payload: errResp,
			
		}
		resp, err := json.Marshal(wsMessage)
		if err != nil {
			logger.ZapLogger.Errorw("JSON Error", "Error", err)
			return
		}
		wrChan <- resp
		return
	}

	wrChan <- resp
}

// func CreateNewHub(logger *zap.SugaredLogger, db db.DbOperations) *NewHub {
// 	return &NewHub{
// 		Id: "hub123",
// 		logger: logger,
// 		ReadMsg: make(chan models.WSMessage, 2),
// 		WriteMsg: make(chan models.WSMessage, 2),
// 		// ClientMap: map[string]*service.Client{},
// 		// Register: make(chan *service.Client),
// 		// Deregister: make(chan *service.Client),
// 		db: db,
// 	}
// }



// func (hub *NewHub) Event(message any) {
// 	logger.ZapLogger.Infoln("Start ProcessEvent")

// 	wsMsg := message.(models.WSMessage)
// 	hub.logger.Infow("ReadMsg", "Type", wsMsg.Type, "Message", string(wsMsg.Payload))

// 	// rollDice := gameplay.CreateRollDice(hub.db, hub.ReadMsg)
// 	// _, err := rollDice.Validate(wsMsg.Payload)
// 	// if err != nil {
// 	// 	return
// 	// }





// }
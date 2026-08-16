package handler

import (
	db "Monopoly/DB"
	gameplay "Monopoly/Gameplay"
	gameroom "Monopoly/Gameroom"
	models "Monopoly/Models"
	"Monopoly/logger"
	"Monopoly/router"
	"encoding/json"

	"go.uber.org/zap"
)

type NewHub struct {
	Id       string
	logger   *zap.SugaredLogger
	ReadMsg  chan models.WSMessage
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
	game     gameplay.Game
	gameRoom *gameroom.GameRoom
}

func CreateHubContoller(g gameplay.Game, gr *gameroom.GameRoom) *GameHubController {
	return &GameHubController{
		game:     g,
		gameRoom: gr,
	}
}

func (hc *GameHubController) HandleHub(req router.Request, readChan chan []byte) {

	// wsMsg := message.(models.WSMessage)
	// hub.Infow("ReadMsg", "Type", wsMsg.Type, "Message", string(wsMsg.Payload))

	action := req.Action
	msg := req.Msg
	reqParam := req.Param

	// Send the Ws Message and get the payload
	payload, err := hc.game.Validate(msg)
	if err != nil {
		logger.ZapLogger.Errorw("Validation Error", "Error", err)
		errResponse := models.RespError{
			Stage:  models.VALIDATION,
			Error:  err.Error(),
			Status: models.FAILED,
		}
		errResp, err := json.Marshal(errResponse)
		if err != nil {
			logger.ZapLogger.Errorw("JSON Error", "Error", err)
			return
		}
		wsMessage := models.WSMessage{
			Type:    action,
			Payload: errResp,
		}
		resp, err := json.Marshal(wsMessage)
		if err != nil {
			logger.ZapLogger.Errorw("JSON Error", "Error", err)
			return
		}
		readChan <- resp
		return
	}

	targetMap, err := hc.game.Play(payload, reqParam)
	if err != nil {
		logger.ZapLogger.Errorw("Play Error", "Error", err)
		errResponse := models.RespError{
			Stage:  models.PLAY,
			Error:  err.Error(),
			Status: models.FAILED,
		}
		errResp, err := json.Marshal(errResponse)
		if err != nil {
			logger.ZapLogger.Errorw("JSON Error", "Error", err)
			return
		}
		wsMessage := models.WSMessage{
			Type:    action,
			Payload: errResp,
		}
		resp, err := json.Marshal(wsMessage)
		if err != nil {
			logger.ZapLogger.Errorw("JSON Error", "Error", err)
			return
		}
		readChan <- resp
		return
	}

	err = hc.game.Response(targetMap, reqParam, readChan)
	if err != nil {

	}

}

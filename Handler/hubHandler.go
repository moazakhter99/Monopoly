package handler

import (
	gameplay "Monopoly/Gameplay"
	models "Monopoly/Models"
	"Monopoly/logger"
	"Monopoly/router"
	"encoding/json"

)

type HubController interface {
	HandleHub(msg []byte, wrChan chan []byte)
}

type GameHubController struct {
	game     gameplay.Game
}

func CreateHubContoller(g gameplay.Game) *GameHubController {
	return &GameHubController{
		game:     g,
	}
}

func (hc *GameHubController) HandleHub(req router.Request, readChan chan []byte) {

	action := req.Action
	msg := req.Msg
	reqParam := req.Param

	// Send the Ws Message and get the payload
	payload, err := hc.game.Validate(msg)
	if err != nil {
		// Either find a different way to send error or use the readChan as error Chan or just rename it
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

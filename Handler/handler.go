package handler

import (
	service "Monopoly/Service"
	"Monopoly/logger"
	"io"
	"net/http"
)

type Controller interface {
	GameHandler(w http.ResponseWriter, r *http.Request)
}

type HandleGameController struct {
	Processor service.Processor
}

func NewGameController(processor service.Processor) *HandleGameController {
	return &HandleGameController{
		Processor: processor,
	}
}

func (game *HandleGameController) GameHandler(w http.ResponseWriter, r *http.Request) {
	logger.ZapLogger.Infoln("Enter Game Handler")

	logger.ZapLogger.Infow("Request URI", "URI", r.RequestURI)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.ZapLogger.Errorw("Request Body", "Error", err)
		return
	}

	req, err := game.Processor.Validate(body)
	if err != nil {
		logger.ZapLogger.Errorw("Request Validation", "Error", err)
	}

	resp, err := game.Processor.ProcessMsg(req)
	if err != nil {
		logger.ZapLogger.Errorw("Request Proccessing", "Error", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(resp)
	logger.ZapLogger.Infoln("Exit Game Handler")
}

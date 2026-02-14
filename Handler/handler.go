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
	Processor service.RequestProcessor
}

func NewGameController(processor service.RequestProcessor) *HandleGameController {
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

	var resp []byte
	uri := r.URL.Path	

	switch uri {
	case "/health":
		resp, err = game.Processor.Health()

	case "/game/create":
		resp, err = game.Processor.CreateGame(body)

	case "/game/join":
		resp, err = game.Processor.JoinGame(body)

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(resp)
	logger.ZapLogger.Infoln("Exit Game Handler")
}

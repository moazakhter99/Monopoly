package handler

import (
	service "Monopoly/Service"
	"Monopoly/logger"

	// "io"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

type HandleWsGameController struct {
	// Processor service.RequestProcessor
	// client *client.Client
	gameHub service.GameHubProcessor

}

func NewWsGameController(gameHub service.GameHubProcessor) *HandleWsGameController {
	return &HandleWsGameController{
		gameHub: gameHub,

	}

}


func (game *HandleWsGameController) WSHandler(w http.ResponseWriter, r *http.Request) {
	logger.ZapLogger.Infoln("Enter Game WebSockert handler")

	if r.Header.Get("Upgrade") != "websocket" {
		http.Error(w, "Not a websocket handshake", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.ZapLogger.Errorw("WebSockert Upgrade", "'Error", err)
		return
	}
	defer conn.Close()

	playerId := ""
	gameId := ""
	gameLog := logger.ZapLogger.With(
		"Player", playerId,
		"GameId", gameId,
	)

	client := service.CreateNewClient(playerId, conn, gameLog)

	go client.ReadMessage()
	go client.WriteMessage()

	// go game.gameHub.ProcessEvent("")


}
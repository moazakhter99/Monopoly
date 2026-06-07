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
	// gameHub service.ClinetProcessor
	gameHub service.GameHubProcessor
	hub *service.GameHub

}

func NewWsGameController(hub *service.GameHub) *HandleWsGameController {
	return &HandleWsGameController{
		hub: hub,

	}

}


func (game *HandleWsGameController) WSHandler(w http.ResponseWriter, r *http.Request) {
	logger.ZapLogger.Infoln("Enter Game WebSockert handler")

	if r.Header.Get("Upgrade") != "websocket" {
		http.Error(w, "Not a websocket handshake", http.StatusBadRequest)
		return
	}
	logger.ZapLogger.Infoln("level 1")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.ZapLogger.Errorw("WebSockert Upgrade", "'Error", err)
		return
	}
	logger.ZapLogger.Infoln("level 2")
	// defer conn.Close()

	// playerId := r.Context().Value("playerId").(string)
    gameId := r.URL.Query().Get("gameId")
    playerId := r.URL.Query().Get("playerId")
	gameLog := logger.ZapLogger.With(
		"Player", playerId,
		"GameId", gameId,
	)
	logger.ZapLogger.Infoln("level 3")

	// game.gameHub.ProcessEvent("")

	// game.gameHub.UpgradeClinet(playerId, conn, gameLog)
	client := service.CreateNewClient(playerId, conn, gameLog, game.hub)
	logger.ZapLogger.Infoln("level 4")
	game.hub.Register <- client
	logger.ZapLogger.Infow("Player Created", "playerId", playerId, "GameId", gameId)


	go client.ReadMessage()
	go client.WriteMessage()

	// go game.gameHub.ProcessEvent("")

	logger.ZapLogger.Infoln("Exit Game WebSocket handler")
}
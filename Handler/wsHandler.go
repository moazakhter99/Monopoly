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
	Processor service.RequestProcessor

}

func NewWsGameController(proccesor service.RequestProcessor) *HandleWsGameController {
	return &HandleWsGameController{
		Processor: proccesor,

	}

}


func (game *HandleWsGameController) GameHandler(w http.ResponseWriter, r *http.Request) {
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



}
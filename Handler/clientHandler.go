package handler

import (
	gameroom "Monopoly/Gameroom"
	"Monopoly/logger"
	"Monopoly/router"

	// "Monopoly/routes"
	"net/http"

	"github.com/gorilla/websocket"
	// "go.uber.org/zap"
)


var cUpgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type HandleClientController struct {
	room 		gameroom.Room
	r			router.Router

}

func CreateClientController(r gameroom.Room, ro router.Router) *HandleClientController {
	return &HandleClientController{
		room: r,
		r: ro,
	}
}


func (cl *HandleClientController) ClientHandler(w http.ResponseWriter, r *http.Request) {
	logger.ZapLogger.Infoln("Enter Client Handler")

	if r.Header.Get("Upgrade") != "websocket" {
		http.Error(w, "Not a websocket handshake", http.StatusBadRequest)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.ZapLogger.Errorw("WebSockert Upgrade", "'Error", err)
		return
	}

    gameId := r.URL.Query().Get("gameId")
    playerId := r.URL.Query().Get("playerId")
	logger.ZapLogger.Infow("Request", "gameId", gameId, "playerId", playerId, "URI", r.RequestURI)

	client := gameroom.CreateNewOtherClinet(cl.r, conn, cl.room, playerId, gameId)

	go client.ReadMessage()
	go client.WriteMessage()

	logger.ZapLogger.Infoln("Exit Client Handler")
}


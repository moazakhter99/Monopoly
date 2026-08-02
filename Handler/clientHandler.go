package handler

import (
	service "Monopoly/Service"
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
	clientProc service.ClinetProcessor

}

func CreateClientController(proc service.ClinetProcessor) *HandleClientController {
	return &HandleClientController{
		clientProc: proc,
	}
}


func (cl *HandleClientController) ClientHandler(w http.ResponseWriter, r *http.Request) {


	// hub := routes.CreateHubContoller()

	// Get gameId
	// Use it to find gameHub
	// Get PlayerId
	// Register Player to that game

	go cl.clientProc.ReadMessage()
	go cl.clientProc.WriteMessage()

}


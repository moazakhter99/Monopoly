package routes

import (
	"Monopoly/DB"
	handler "Monopoly/Handler"
	"Monopoly/Service/gameService"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)




var GameSubRouter = func (router *mux.Router, monopolyDB db.DbOperations, logger *zap.SugaredLogger) {

	createGameReq := service.CreateGameReq(monopolyDB, logger)
	gameReqHandler := handler.NewGameController(createGameReq)
	router.HandleFunc("/create", gameReqHandler.GameHandler).Methods("POST")

}
package routes

import (
	"Monopoly/DB"
	handler "Monopoly/Handler"
	"Monopoly/Service"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)


var GameSubRouter = func (router *mux.Router, monopolyDB db.DbOperations, logger *zap.SugaredLogger) {

	createGameReq := service.CreateGameReq(monopolyDB, logger)
	createGameHandler := handler.NewGameController(createGameReq)
	router.HandleFunc("/create", createGameHandler.GameHandler).Methods("POST")

	joinGameReq := service.CreateJoinGameReq(monopolyDB, logger)
	joinGameHandler := handler.NewGameController(joinGameReq)
	router.HandleFunc("/join", joinGameHandler.GameHandler).Methods("POST")

}
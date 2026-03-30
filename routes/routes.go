package routes

import (
	db "Monopoly/DB"
	handler "Monopoly/Handler"
	"Monopoly/Service"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)


var InitGameSubRouter = func (router *mux.Router, db db.DbOperations, logger *zap.SugaredLogger) {


	createGameReq := service.CreateGameReq(db, logger)
	createGameReqHandler := handler.NewGameController(createGameReq)
	router.HandleFunc("/create", createGameReqHandler.GameHandler).Methods("POST")


	joinGameReq := service.CreateJoinGameReq(db, logger)
	joinGameHandler := handler.NewGameController(joinGameReq)
	router.HandleFunc("/join", joinGameHandler.GameHandler).Methods("POST")

	startGameReq := service.CreateStartGameReq(db, logger)
	startGameHandler := handler.NewGameController(startGameReq)
	router.HandleFunc("/start", startGameHandler.GameHandler).Methods("POST")

}

// var GameSubRouter = func (router *mux.Router, proc service.RequestProcessor) {


// }

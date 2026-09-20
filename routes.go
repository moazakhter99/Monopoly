package main

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

var GetInfoRouter = func (router *mux.Router, db db.DbOperations) {

	createJailReq := service.CreateJailReq(db)
	jailReqHandler := handler.NewGameController(createJailReq)
	router.HandleFunc("/jail", jailReqHandler.GameHandler)

	createCommunityChestReq := service.CreateCommunityChestReq(db)
	communityChestHandler := handler.NewGameController(createCommunityChestReq)
	router.HandleFunc("/communityChest", communityChestHandler.GameHandler)

	createChanceReq := service.CreateChanceReq(db)
	chanceHandler := handler.NewGameController(createChanceReq)
	router.HandleFunc("/chance", chanceHandler.GameHandler)

}

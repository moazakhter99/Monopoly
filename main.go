package main

import (
	db "Monopoly/DB"
	"Monopoly/DB/postgres"
	"Monopoly/DB/sqlLite"
	gamehub "Monopoly/Gamehub"
	gameroom "Monopoly/Gameroom"
	handler "Monopoly/Handler"
	service "Monopoly/Service"
	"Monopoly/load"
	"Monopoly/logger"
	"Monopoly/router"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {
	load.Env()
	logger.Logger()

	port := viper.GetString("PORT")
	database := viper.GetString("DATABASE")
	chanBufferSize := viper.GetInt("BUFFER_SIZE")
	logger.ZapLogger.Infow("Game Running on", "PORT", port)

	r := mux.NewRouter()
	wsRouter := router.NewRouter(chanBufferSize)

	var MonopolyDB db.DbOperations

	switch database {
	case "POSTGRES":
		postgres, err := postgres.OpenDatabase()
		if err != nil {
			logger.ZapLogger.Panic("Database Connection", "Error", err)
			return
		}
		MonopolyDB = postgres

	case "SQLLITE":
		sqlLite, err := sqlLite.OpenDatabase()
		if err != nil {
			logger.ZapLogger.Panic("Database Connection", "Error", err)
			return
		}
		MonopolyDB = sqlLite

	default:
		logger.ZapLogger.Panic("Unknown Database", "Database", database)
		return

	}


	healthReq := service.CreateHealthReq(MonopolyDB)
	healthHandler := handler.NewGameController(healthReq)
	r.HandleFunc("/health", healthHandler.GameHandler).Methods("GET")

	initGameRouter := r.PathPrefix("/initGame").Subrouter()
	InitGameSubRouter(initGameRouter, MonopolyDB, logger.ZapLogger)

	getInfoRouter := r.PathPrefix("/info").Subrouter().Methods("GET")
	GetInfoRouter(getInfoRouter.Subrouter(), MonopolyDB)

	gameHub := service.CreateNewGameHub(logger.ZapLogger, MonopolyDB)
	wsHandler := handler.NewWsGameController(gameHub)
	r.HandleFunc("/oldws", wsHandler.WSHandler)
	// go run(gameHub)

	client := gameroom.CreateOtherClinet(wsRouter)
	gameRoom := gameroom.CreateGameRoom()
	clientHandler := handler.CreateClientController(client, gameRoom)
	r.HandleFunc("/ws", clientHandler.ClientHandler)
	go gamehub.Monopoly(MonopolyDB, wsRouter)

	http.ListenAndServe(":"+port, r)

}


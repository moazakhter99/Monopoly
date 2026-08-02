package main

import (
	db "Monopoly/DB"
	"Monopoly/DB/postgres"
	"Monopoly/DB/sqlLite"
	gamehub "Monopoly/Gamehub"
	handler "Monopoly/Handler"
	"Monopoly/router"
	service "Monopoly/Service"
	"Monopoly/load"
	"Monopoly/logger"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {
	load.Env()
	logger.Logger()

	port := viper.GetString("PORT")
	database := viper.GetString("DATABASE")
	logger.ZapLogger.Infow("Game Running on", "PORT", port)

	r := mux.NewRouter()
	wsRouter := router.NewRouter()

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

	client := service.CreateOtherClinet(wsRouter)
	clientHandler := handler.CreateClientController(client)
	r.HandleFunc("/ws", clientHandler.ClientHandler)
	go gamehub.Monopoly(MonopolyDB, wsRouter)

	http.ListenAndServe(":"+port, r)

}


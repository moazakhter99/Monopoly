package main

import (
	db "Monopoly/DB"
	"Monopoly/DB/postgres"
	"Monopoly/DB/sqlLite"
	handler "Monopoly/Handler"
	service "Monopoly/Service"
	"Monopoly/load"
	"Monopoly/logger"
	"Monopoly/routes"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {
	load.Env()
	logger.Logger()

	port := viper.GetString("PORT")
	logger.ZapLogger.Infow("Game Running on", "PORT", port)
	router := mux.NewRouter()

	var MonopolyDB db.DbOperations
	database := viper.GetString("DATABASE")

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

	reqProc := service.CreateNewRequestProcessor(MonopolyDB, logger.ZapLogger)

	healthHandler := handler.NewGameController(reqProc)
	router.HandleFunc("/health", healthHandler.GameHandler).Methods("GET")

	gameRouter := router.PathPrefix("/game").Subrouter()
	routes.GameSubRouter(gameRouter, reqProc)

	http.ListenAndServe(":"+port, router)

}

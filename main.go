package main

import (
	db "Monopoly/DB"
	"Monopoly/DB/postgres"
	"Monopoly/DB/sqlLite"
	"Monopoly/load"
	"Monopoly/logger"

	"github.com/spf13/viper"
)


func main() {
	load.Env()
	logger.Logger()

	var MonopolyDB db.DbOperations
	logger.ZapLogger.Infoln("Start Database")
	database := viper.GetString("DATABASE")

	switch database {
	case "POSTGRES":
		postgres, err := postgres.OpenDatabase()
		if err != nil {
			logger.ZapLogger.Panicln("Database Connection", "Error", err)
			return
		}
		MonopolyDB = postgres

	case "SQLLITE":
		sqlLite, err := sqlLite.OpenDatabase()
		if err != nil {
			logger.ZapLogger.Panicln("Database Connection", "Error", err)
			return
		}
		MonopolyDB = sqlLite

	default:		
		logger.ZapLogger.Panicln("Unknown Database", "Database", database)
		return

	}


}
package main

import (
	"Monopoly/load"
	"Monopoly/logger"
)


func main() {
	load.Env()
	logger.Logger()

	logger.ZapLogger.Infoln("Start Database")

}
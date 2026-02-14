package routes

import (
	handler "Monopoly/Handler"
	"Monopoly/Service"

	"github.com/gorilla/mux"
)


var GameSubRouter = func (router *mux.Router, proc service.RequestProcessor) {

	createGameHandler := handler.NewGameController(proc)
	router.HandleFunc("/create", createGameHandler.GameHandler).Methods("POST")

	joinGameHandler := handler.NewGameController(proc)
	router.HandleFunc("/join", joinGameHandler.GameHandler).Methods("POST")

}
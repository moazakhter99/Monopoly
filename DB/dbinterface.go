package db

import models "Monopoly/Models"


type DbOperations interface {
	Ping() (err error)	
	InsertGame(gameId, matchId string) (err error)
	InsertPlayer(player *models.Player, gameId string) (err error)
	GetGameFromMatchId(matchId string) (gameId string, err error)

}
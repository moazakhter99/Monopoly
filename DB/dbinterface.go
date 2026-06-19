package db

import models "Monopoly/Models"


type DbOperations interface {
	Ping() (err error)	
	InsertGame(gameId, matchId string) (err error)
	InsertPlayer(player *models.Player, gameId string) (err error)
	GetGameFromMatchId(matchId string) (gameId string, err error)
	GetPlayerInfoById(playerId string) (player *models.Player, err error)
	GetBlockState(position int, gameId string) (state bool, blockId string, err error)
	UpdatePlayerPos(playerId string, position int) (err error)
	GetBlockInfoById(blockId string) (block *models.Block, err error)
	UpdatePlayerCard(playerId, gameId, blockId string) (err error)
	GetCardAction(cardNo string) (action string, err error)
	GetPlayerCash(playerId string) (cash int, err error)
	UpdatePlayerCash(playerId string, cash int) (err error)

}
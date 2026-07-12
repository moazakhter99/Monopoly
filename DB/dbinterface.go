package db

import models "Monopoly/Models"


type DbOperations interface {
	Ping() (err error)	
	InsertGame(gameId, matchId string) (err error)
	InsertPlayer(player *models.Player, gameId string) (err error)
	GetGameFromMatchId(matchId string) (gameId string, err error)
	GetPlayerInfoById(playerId string) (player *models.Player, err error)
	GetBlockState(position int, gameId string) (block *models.Block, err error)
	UpdatePlayerPos(playerId string, position int) (err error)
	GetBlockInfoById(blockId string) (block *models.Block, err error)
	UpdatePlayerCard(playerId, gameId, blockId, status string) (err error)
	GetCardAction(cardNo string) (action string, err error)
	GetPlayerCash(playerId string) (cash int, err error)
	UpdatePlayerCash(playerId string, cash int) (err error)
	GetPosByBlockName(blockName string) (pos int, err error)
	GetBlockInfoByBlockType(blockType string) (jailInfo []models.Jail, err error)
	GetCardOwnership(blockId, gameId string) (playerId string, err error)
	UpdateGetOutOfJailCard(playerId, gameId string) (err error)
	DeleteGetOutOfJailCard(playerId, GameId, cardId string) (err error)
	UpdatePlayerStatus(playerId, count string) (err error)
	GetNextPlayer(gameId string, seq int) (playerId string, err error)
	GetPlayerSeqAndCount(playerId string) (seq, count int, err error)
	GetBlockOwner(blockID, gameID string) (playerID string, err error)
	GetCardOwnershipStatus(playerId, blockId string) (status string, err error)
	GetCardOwnerCount(playerId, gameId string) (count int, err error)
	GetPlayerStatusList(playerId, gameId string) (statusList []string, err error)

}
package gameplay

import (
	db "Monopoly/DB"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
)

func RollDice(request json.RawMessage, readCh chan<- models.WSMessage) (response json.RawMessage) {
	logger.ZapLogger.Infoln("Roll Dice")
	var req models.Request

	err := json.Unmarshal(request, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.ROLLDICE, "JSON Error", err)
		return
	}
	playerId := req.PlayerId

	diceVal := diceRoll()

	resp := models.RespDiceRoll{
		PlayerId: playerId,
		GameId:   req.GameId,
		DiceVal:  diceVal,
	}

	response, err = json.Marshal(resp)
	if err != nil {
		logger.ZapLogger.Errorw(models.ROLLDICE, "JSON Error", err)
		return
	}

	go func() {
		newPosReq := models.ReqMovePos{
			PlayerId: req.PlayerId,
			GameId:   req.GameId,
			UpdateBy: diceVal,
		}

		newPosPayload, err := json.Marshal(newPosReq)
		if err != nil {
			logger.ZapLogger.Errorw(models.ROLLDICE, "JSON Error", err)
			return
		}

		movePos := models.WSMessage{
			Type:    models.MOVEPOS,
			Payload: newPosPayload,
		}

		readCh <- movePos

	}()

	return
}

func MovePos(request json.RawMessage, db db.DbOperations) (response json.RawMessage) {
	logger.ZapLogger.Info("Move Position")
	var req models.ReqMovePos

	err := json.Unmarshal(request, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.MOVEPOS, "JSON Error", err)
		return
	}

	player, err := db.GetPlayerInfoById(req.PlayerId)
	if err != nil {
		logger.ZapLogger.Errorw(models.MOVEPOS, "DB Error", err)
		return
	}
	logger.ZapLogger.Infow(models.MOVEPOS, "Current Position", player.Pos)

	newPos := updatePos(req.UpdateBy, player.Pos)
	err = db.UpdatePlayerPos(req.PlayerId, newPos)
	if err != nil {
		logger.ZapLogger.Errorw(models.MOVEPOS, "DB Error", err)
	}

	state, blockId, err := db.GetBlockState(newPos, req.GameId)
	if err != nil {
		logger.ZapLogger.Errorw(models.MOVEPOS, "DB Error", err)
		return
	}
	logger.ZapLogger.Infow(models.MOVEPOS, "Block State", state, "New Position", newPos, "Block Id", blockId)

	resp := models.RespMovePos{
		PlayerId: req.PlayerId,
		GameId:   req.GameId,
		BlockId:  blockId,
		NewPos:   newPos,
		Sold:     state,
	}

	response, err = json.Marshal(resp)
	if err != nil {
		logger.ZapLogger.Errorw(models.MOVEPOS, "JSON Error", err)
		return
	}

	return
}

func BuyBlock(request json.RawMessage, db db.DbOperations) (response json.RawMessage) {
	logger.ZapLogger.Infoln("Buy Block")
	var req models.ReqBuyBlock
	var updatedCash int
	var buy bool

	err := json.Unmarshal(request, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.BUYBLOCK, "JSON Error", err)
		return
	}

	cash, err := db.GetPlayerCash(req.PlayerId)
	if err != nil {
		logger.ZapLogger.Errorw(models.BUYBLOCK, "DB Error", err)
	}

	if cash < req.Price {
		logger.ZapLogger.Infow("Player Low on Cash", "Player Id", req.PlayerId, "Cash", cash)
		buy = false

	} else {
		updatedCash = cash - req.Price
		buy = true

	}

	err = db.UpdatePlayerCard(req.PlayerId, req.GameId, req.BlockId)
	if err != nil {
		logger.ZapLogger.Errorw(models.BUYBLOCK, "DB Error", err)
	}

	logger.ZapLogger.Infow(models.BUYBLOCK, "Game Id", req.GameId, "Player Id", req.PlayerId, "Block Id", req.BlockId)
	resp := models.RespBuyBlock{
		PlayerId: req.PlayerId,
		GameId:   req.GameId,
		BlockId:  req.BlockId,
		Buy:      buy,
		Cash:     updatedCash,
	}

	err = db.UpdatePlayerCash(req.PlayerId, updatedCash)
	if err != nil {
		logger.ZapLogger.Errorw(models.BUYBLOCK, "DB Error", err)
	}

	response, err = json.Marshal(resp)
	if err != nil {
		logger.ZapLogger.Errorw(models.BUYBLOCK, "JSON Error", err)
		return
	}

	return
}
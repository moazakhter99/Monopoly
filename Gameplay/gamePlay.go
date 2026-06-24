package gameplay

import (
	db "Monopoly/DB"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
	"strconv"
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

	block, err := db.GetBlockState(newPos, req.GameId)
	if err != nil {
		logger.ZapLogger.Errorw(models.MOVEPOS, "DB Error", err)
		return
	}
	logger.ZapLogger.Infow(models.MOVEPOS, "Block State", block.State, "New Position", newPos, "Block Id", block.BlockId)

	resp := models.RespMovePos{
		PlayerId: req.PlayerId,
		GameId:   req.GameId,
		BlockId:  block.BlockId,
		NewPos:   newPos,
		Sold:     block.State,
		Type:     block.Type,
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

func ActionCard(request json.RawMessage, db db.DbOperations, readCh chan<- models.WSMessage) (response json.RawMessage) {
	logger.ZapLogger.Infoln("Enter Action Block")
	var req models.ReqActionCard

	err := json.Unmarshal(request, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.ACTIONCARD, "JSON Error", err)
		return
	}

	resp := models.RespActionCard{
		PlayerId: req.PlayerId,
		GameId:   req.GameId,
	}

	switch req.Type {

	case models.COMMUNITYCHEST:
		action, err := db.GetCardAction(req.CardId)
		if err != nil {
			logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
			return
		}

		chestValue, err := strconv.Atoi(action)
		if err != nil {
			logger.ZapLogger.Errorw(models.ACTIONCARD, "Conversion Error", err)
			return
		}
		resp.Cash = req.Cash + chestValue

	case models.CHANCE:
		action, err := db.GetCardAction(req.CardId)
		if err != nil {
			logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
			return
		}

		switch action {

		case models.JUMPTOMUMBAI:
			resp.Pos, err = db.GetPosByBlockName("Mumbai")
			if err != nil {
				logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
				return
			}

		case models.GOTOSTART:
			resp.Pos, err = db.GetPosByBlockName("Go Start")
			if err != nil {
				logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
				return
			}

		case models.HOUSEHOTELFINE:
			return

		case models.GETOUTOFJAIL:
			err = db.UpdateGetOutOfJailCard(req.PlayerId, req.GameId)
			if err != nil {
				logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
				return
			}

		default:
			value, err := strconv.Atoi(action)
			if err != nil {
				logger.ZapLogger.Infow(models.ACTIONCARD, "Conversion Error", err, "Invalid Action", action)
				return
			}
			resp.Cash = req.Cash + value

		}

	case models.INCOMETAX:
		resp.Cash = req.Cash - req.Price

	case models.JAIL:
		resp.Pos = req.Pos

		go func() {
			jailInfo, err := db.GetBlockInfoByBlockType(models.JAIL)
			if err != nil {
				logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
				return
			}

			ownerId, err := db.GetCardOwnership(req.BlockId, req.GameId)
			if err != nil {
				logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
				return
			}

			reqjail := models.RespJail{
				PlayerId:   req.PlayerId,
				GameId:     req.GameId,
				Jail:       jailInfo,
				GetOutCard: ownershipConfirm(ownerId, req.PlayerId),
				InJail:     true,
			}

			jailReq, err := json.Marshal(reqjail)
			if err != nil {
				logger.ZapLogger.Errorw(models.ROLLDICE, "JSON Error", err)
				return
			}

			jailWsMsg := models.WSMessage{
				Type:    models.JAIL,
				Payload: jailReq,
			}

			readCh <- jailWsMsg

		}()

	case models.FREEPARKING:

	case models.GOTOJAIL:

		go func() {
			jailInfo, err := db.GetBlockInfoByBlockType(models.JAIL)
			if err != nil {
				logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
				return
			}

			jailPos, err := db.GetPosByBlockName("Jail")
			if err != nil {
				logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
				return
			}

			reqjail := models.RespJail{
				PlayerId: req.PlayerId,
				GameId:   req.GameId,
				Jail:     jailInfo,
				NewPos:   jailPos,
				// Get this from DB
				GetOutCard: false,
				InJail:     true,
			}

			jailReq, err := json.Marshal(reqjail)
			if err != nil {
				logger.ZapLogger.Errorw(models.ROLLDICE, "JSON Error", err)
				return
			}

			jailWsMsg := models.WSMessage{
				Type:    models.JAIL,
				Payload: jailReq,
			}

			readCh <- jailWsMsg

		}()

	case models.PROPERTYTAX:
		resp.Cash = req.Cash - req.Price

	case models.GOTOSTART:
		resp.Cash = req.Cash - req.Price

	}

	response, err = json.Marshal(resp)
	if err != nil {
		logger.ZapLogger.Errorw(models.ROLLDICE, "JSON Error", err)
		return
	}

	return
}

func Jail(request json.RawMessage, db db.DbOperations) (response json.RawMessage) {
	logger.ZapLogger.Infoln("Enter Jail")
	var req models.ReqJail
	var updatedCash int
	var inJail bool

	err := json.Unmarshal(request, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.JAIL, "JSON Error", err)
		return
	}

	switch req.JailId {

	case "Jail0":
		updatedCash = req.Cash - 500
		err = db.UpdatePlayerCash(req.PlayerId, updatedCash)
		if err != nil {
			logger.ZapLogger.Errorw(models.JAIL, "DB Error", err)
			return
		}

		inJail = false

	case "Jail1":
		err = db.DeleteGetOutOfJailCard(req.PlayerId, req.GameId)
		if err != nil {
			logger.ZapLogger.Errorw(models.JAIL, "DB Error", err)
			return
		}
		inJail = false
		updatedCash = req.Cash

	case "Jail2":
		err = db.UpdatePlayerStatus(req.PlayerId, "3")
		if err != nil {
			logger.ZapLogger.Errorw(models.JAIL, "DB Error", err)
			return
		}
		inJail = true
		updatedCash = req.Cash

	}
	
	resp := models.RespJail{
		PlayerId: req.PlayerId,
		GameId: req.GameId,
		Cash: updatedCash,
		InJail: inJail,
	}

	response, err = json.Marshal(resp)
	if err != nil {
		logger.ZapLogger.Errorw(models.JAIL, "JSON Error", err)
		return
	}

	logger.ZapLogger.Infoln("Exit Jail")
	return
}
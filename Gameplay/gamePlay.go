package gameplay

import (
	db "Monopoly/DB"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
	"strconv"
)

func RollDice(request json.RawMessage, client *models.Client, readCh chan<- models.WSMessage) (response json.RawMessage) {
	logger.ZapLogger.Infoln("Roll Dice")
	var req models.Request

	err := json.Unmarshal(request, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.ROLLDICE, "JSON Error", err)
		return
	}

	diceVal := diceRoll()

	resp := models.RespDiceRoll{
		DiceVal:  diceVal,
	}

	response, err = json.Marshal(resp)
	if err != nil {
		logger.ZapLogger.Errorw(models.ROLLDICE, "JSON Error", err)
		return
	}

	go func() {
		newPosReq := models.ReqMovePos{
			UpdateBy: diceVal,
		}

		newPosPayload, err := json.Marshal(newPosReq)
		if err != nil {
			logger.ZapLogger.Errorw(models.ROLLDICE, "JSON Error", err)
			return
		}

		movePos := models.WSMessage{
			Type:    models.MOVEPOS,
			Client:  client,
			Payload: newPosPayload,
		}

		readCh <- movePos

	}()

	return
}

func MovePos(request json.RawMessage, client *models.Client, db db.DbOperations) (response json.RawMessage) {
	logger.ZapLogger.Info("Move Position")
	var req models.ReqMovePos

	err := json.Unmarshal(request, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.MOVEPOS, "JSON Error", err)
		return
	}
	playerId := client.PlayerId
	gameId := client.GameId

	player, err := db.GetPlayerInfoById(playerId)
	if err != nil {
		logger.ZapLogger.Errorw(models.MOVEPOS, "DB Error", err)
		return
	}
	logger.ZapLogger.Infow(models.MOVEPOS, "Current Position", player.Pos)

	newPos := updatePos(req.UpdateBy, player.Pos)
	err = db.UpdatePlayerPos(playerId, newPos)
	if err != nil {
		logger.ZapLogger.Errorw(models.MOVEPOS, "DB Error", err)
	}

	block, err := db.GetBlockState(newPos, gameId)
	if err != nil {
		logger.ZapLogger.Errorw(models.MOVEPOS, "DB Error", err)
		return
	}
	logger.ZapLogger.Infow(models.MOVEPOS, 
		"Block State", block.State, 
		"New Position", newPos, 
		"Block Id", block.BlockId)

	resp := models.RespMovePos{
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

func BuyBlock(request json.RawMessage, client *models.Client, db db.DbOperations, readCh chan<- models.WSMessage) (response json.RawMessage) {
	logger.ZapLogger.Infoln("Buy Block")
	var req models.ReqBuyBlock
	var updatedCash int
	var buy bool

	err := json.Unmarshal(request, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.BUYBLOCK, "JSON Error", err)
		return
	}
	playerId := client.PlayerId
	gameId := client.GameId

	cash, err := db.GetPlayerCash(playerId)
	if err != nil {
		logger.ZapLogger.Errorw(models.BUYBLOCK, "DB Error", err)
	}

	if cash < req.Price {
		logger.ZapLogger.Infow("Player Low on Cash", "Player Id", playerId, "Cash", cash)
		buy = false

	} else {
		updatedCash = cash - req.Price
		buy = true

	}

	err = db.UpdatePlayerCard(playerId, gameId, req.BlockId)
	if err != nil {
		logger.ZapLogger.Errorw(models.BUYBLOCK, "DB Error", err)
	}

	logger.ZapLogger.Infow(models.BUYBLOCK, "Game Id", gameId, "Player Id", playerId, "Block Id", req.BlockId)

	go func() {

		changePlayerReq := models.Request{

		}

		payload, err := json.Marshal(changePlayerReq)
		if err != nil {
			logger.ZapLogger.Errorw(models.CHANGEPLAYER, "JSON Error", err)
			return
		}

		changePlayer := models.WSMessage{
			Type: models.CHANGEPLAYER,
			Client: client,
			Payload: payload,
		}

		readCh <- changePlayer

	}()

	resp := models.RespBuyBlock{
		BlockId:  req.BlockId,
		Buy:      buy,
		Cash:     updatedCash,
		ChangePlayer: true,
	}

	err = db.UpdatePlayerCash(playerId, updatedCash)
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

func ActionCard(request json.RawMessage, client *models.Client, db db.DbOperations, readCh chan<- models.WSMessage) (response json.RawMessage) {
	logger.ZapLogger.Infoln("Enter Action Block")
	var req models.ReqActionCard

	err := json.Unmarshal(request, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.ACTIONCARD, "JSON Error", err)
		return
	}
	playerId := client.PlayerId
	gameId := client.GameId

	resp := models.RespActionCard{
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
			err = db.UpdateGetOutOfJailCard(playerId, gameId)
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

			ownerId, err := db.GetCardOwnership(req.BlockId, gameId)
			if err != nil {
				logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
				return
			}

			reqjail := models.RespJail{
				Jail:       jailInfo,
				GetOutCard: ownershipConfirm(ownerId, playerId),
				InJail:     true,
			}

			jailReq, err := json.Marshal(reqjail)
			if err != nil {
				logger.ZapLogger.Errorw(models.ROLLDICE, "JSON Error", err)
				return
			}

			jailWsMsg := models.WSMessage{
				Type:    models.JAIL,
				Client:  client,
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
				Client:  client,
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

func Jail(request json.RawMessage, client *models.Client, db db.DbOperations) (response json.RawMessage) {
	logger.ZapLogger.Infoln("Enter Jail")
	var req models.ReqJail
	var updatedCash int
	var inJail bool

	err := json.Unmarshal(request, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.JAIL, "JSON Error", err)
		return
	}
	playerId := client.PlayerId
	gameId := client.GameId

	switch req.JailId {

	case "Jail0":
		updatedCash = req.Cash - 500
		err = db.UpdatePlayerCash(playerId, updatedCash)
		if err != nil {
			logger.ZapLogger.Errorw(models.JAIL, "DB Error", err)
			return
		}

		inJail = false

	case "Jail1":
		err = db.DeleteGetOutOfJailCard(playerId, gameId)
		if err != nil {
			logger.ZapLogger.Errorw(models.JAIL, "DB Error", err)
			return
		}
		inJail = false
		updatedCash = req.Cash

	case "Jail2":
		err = db.UpdatePlayerStatus(playerId, "3")
		if err != nil {
			logger.ZapLogger.Errorw(models.JAIL, "DB Error", err)
			return
		}
		inJail = true
		updatedCash = req.Cash

	}
	
	resp := models.RespJail{
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

func ChangePlayer(request json.RawMessage, client *models.Client, db db.DbOperations) (targetMap map[string][]byte) {
	logger.ZapLogger.Infoln("Enter Change Player")
	var req models.Request
	targetMap = make(map[string][]byte, 2)

	err := json.Unmarshal(request, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.JAIL, "JSON Error", err)
		return
	}
	playerId := client.PlayerId
	gameId := client.GameId

	seq, count, err := db.GetPlayerSeqAndCount(playerId)

	nextPlayerId, err := db.GetNextPlayer(gameId, nextSeq(seq, count))

	nextPlayer := models.RespChangePlayer{
		NextPlayer: nextPlayerId,
		Playing: true,
	}
	nextResp, err := json.Marshal(nextPlayer)
	if err != nil {
		logger.ZapLogger.Errorw(models.BUYBLOCK, "JSON Error", err)
		return
	}
	targetMap[nextPlayerId] = nextResp

	currPlayer := models.RespChangePlayer{
		NextPlayer: nextPlayerId,
		Playing: false,
	}
	currResp, err := json.Marshal(currPlayer)
	if err != nil {
		logger.ZapLogger.Errorw(models.BUYBLOCK, "JSON Error", err)
		return
	}
	targetMap[playerId] = currResp


	logger.ZapLogger.Infoln("Exit Change Player")
	return
}

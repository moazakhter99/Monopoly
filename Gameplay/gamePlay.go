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
	var diceVal int

	err := json.Unmarshal(request, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.ROLLDICE, "JSON Error", err)
		return
	}

	if req.Roll {
		diceVal = diceRoll()

	}

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

func MovePos(request json.RawMessage, client *models.Client, db db.DbOperations, readCh chan<- models.WSMessage) (response json.RawMessage) {
	logger.ZapLogger.Info("Move Position")
	var req models.ReqMovePos
	var resp models.RespMovePos

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

	switch block.OwnerId {
	case "":
		// Buy or Action Card
		if block.Type == models.SPECIALCARD {
			resp = models.RespMovePos{
				BlockId:  block.BlockId,
				NewPos:   newPos,
				Type:     block.Type,
				BlockName: block.BlockName,
			}
		} else {
			resp = models.RespMovePos{
				BlockId:  block.BlockId,
				NewPos:   newPos,
				Type:     block.Type,
			}

		}

	case playerId:
		// Already owned by the player 
		resp = models.RespMovePos{
			BlockId:  block.BlockId,
			NewPos:   newPos,
			State:    models.OWNED,
			Type:     block.Type,
			OwnerId:  playerId,
		}
		logger.ZapLogger.Infow(models.CHANGEPLAYER, "Current Player", playerId)
		go callChangePlayer(client, readCh)

	default:
		// Pay rent
		resp = models.RespMovePos{
			BlockId:  block.BlockId,
			NewPos:   newPos,
			State:    models.SOLD,
			Type:     block.Type,
			OwnerId:  block.OwnerId,
		}

	}

	logger.ZapLogger.Infow(models.MOVEPOS, 
		"Block Sold", block.State, 
		"New Position", newPos, 
		"Block Id", block.BlockId)

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
		return
	}

	if cash < req.Price {
		logger.ZapLogger.Infow("Player Low on Cash", "Player Id", playerId, "Cash", cash)
		buy = false

	} else {
		updatedCash = cash - req.Price
		buy = true

	}

	// Calculate Status for (colour, house, hotel)
	status := models.BASE

	err = db.UpdatePlayerCard(playerId, gameId, req.BlockId, status)
	if err != nil {
		logger.ZapLogger.Errorw(models.BUYBLOCK, "DB Error", err)
		return
	}

	logger.ZapLogger.Infow(models.BUYBLOCK, "Game Id", gameId, "Player Id", playerId, "Block Id", req.BlockId)
	resp := models.RespBuyBlock{
		BlockId:  req.BlockId,
		Buy:      buy,
		Cash:     updatedCash,
		ChangePlayer: true,
	}

	err = db.UpdatePlayerCash(playerId, updatedCash)
	if err != nil {
		logger.ZapLogger.Errorw(models.BUYBLOCK, "DB Error", err)
		return
	}
	go callChangePlayer(client, readCh)

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
		InJail: false,
	}
	// Decide where where to cash value from
	cash, err := db.GetPlayerCash(playerId)
	if err != nil {
		logger.ZapLogger.Errorw(models.BUYBLOCK, "DB Error", err)
		return
	}
	req.Cash = cash
	resp.Cash = cash
	logger.ZapLogger.Infow(models.ACTIONCARD, "Current Cash", req.Cash)

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
		logger.ZapLogger.Infow(models.ACTIONCARD, "Update Cash", chestValue)

	case models.CHANCE:
	
		action, err := db.GetCardAction(req.CardId)
		if err != nil {
			logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
			return
		}

		logger.ZapLogger.Infow(models.ACTIONCARD, "Action", action)
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
			resp.Cash = cash + 2000

		case models.HOUSEHOTELFINE:
			statusList, err := db.GetPlayerStatusList(playerId, gameId)
			if err != nil {
				logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
				return
			}

			fine := houseHotelFine(statusList)
			resp.Cash = req.Cash - fine

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
			logger.ZapLogger.Infow(models.ACTIONCARD, "Value", value)

		}

	case models.INCOMETAX:
		resp.Cash = req.Cash - req.Price

	case models.JAIL:
		resp.Pos = req.Pos
		resp.InJail = true
		status := models.BLOCKED + "_3"

		err = db.UpdatePlayerStatus(playerId, status)
		if err != nil {
			logger.ZapLogger.Errorw(models.JAIL, "DB Error", err)
			return
		}

	case models.FREEPARKING:

	case models.GOTOJAIL:

		jailPos, err := db.GetPosByBlockName("Jail")
		if err != nil {
			logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
			return
		}

		resp.Pos = jailPos
		resp.InJail = true
		status := models.BLOCKED + "_3"

		err = db.UpdatePlayerStatus(playerId, status)
		if err != nil {
			logger.ZapLogger.Errorw(models.JAIL, "DB Error", err)
			return
		}

	case models.PROPERTYTAX:
		resp.Cash = req.Cash - req.Price

	case models.GOTOSTART:
		resp.Cash = req.Cash - req.Price

	default:
		logger.ZapLogger.Errorw(models.ACTIONCARD, "Invalid Type", req.Type)
	}

	// Update Player Remaining

	response, err = json.Marshal(resp)
	if err != nil {
		logger.ZapLogger.Errorw(models.ROLLDICE, "JSON Error", err)
		return
	}

	logger.ZapLogger.Infow(models.ACTIONCARD, "Resp", string(response))
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

	case "Jail1":
		updatedCash = req.Cash - 500
		logger.ZapLogger.Infow(models.JAIL, "Updated Cash", updatedCash)
		err = db.UpdatePlayerCash(playerId, updatedCash)
		if err != nil {
			logger.ZapLogger.Errorw(models.JAIL, "DB Error", err)
			return
		}

		inJail = false

		err = db.UpdatePlayerStatus(playerId, "")
		if err != nil {
			logger.ZapLogger.Errorw(models.JAIL, "DB Error", err)
			return
		}

	case "Jail2":
		err = db.DeleteGetOutOfJailCard(playerId, gameId, "Special Card0")
		if err != nil {
			logger.ZapLogger.Errorw(models.JAIL, "DB Error", err)
			return
		}
		inJail = false
		updatedCash = req.Cash

		err = db.UpdatePlayerStatus(playerId, "")
		if err != nil {
			logger.ZapLogger.Errorw(models.JAIL, "DB Error", err)
			return
		}

	case "Jail3":
		// The player is blocked at the front end and they are also checked at backend
		inJail = req.InJail
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

	logger.ZapLogger.Infow(models.JAIL, "Resp", string(response))
	logger.ZapLogger.Infoln("Exit Jail")
	return
}

func ChangePlayer(request json.RawMessage, client *models.Client, db db.DbOperations) (targetMap map[string][]byte) {
	logger.ZapLogger.Infoln("Enter Change Player")
	var req models.Request
	targetMap = make(map[string][]byte, 2)

	err := json.Unmarshal(request, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.CHANGEPLAYER, "JSON Error", err)
		return
	}
	playerId := client.PlayerId
	gameId := client.GameId

	seq, count, err := db.GetPlayerSeqAndCount(playerId)
	if err != nil {
		logger.ZapLogger.Errorw(models.CHANGEPLAYER, "DB Error", err)
		return
	}
	logger.ZapLogger.Infow(models.CHANGEPLAYER, "CurrSeq", seq, "currCount", count)
	nextPlayerId, err := db.GetNextPlayer(gameId, nextSeq(seq, count))
	if err != nil {
		logger.ZapLogger.Errorw(models.CHANGEPLAYER, "DB Error", err)
		return
	}

	nextPlayer := models.RespChangePlayer{
		NextPlayer: nextPlayerId,
		Playing: true,
	}
	nextResp, err := json.Marshal(nextPlayer)
	if err != nil {
		logger.ZapLogger.Errorw(models.CHANGEPLAYER, "JSON Error", err)
		return
	}
	targetMap[nextPlayerId] = nextResp

	currPlayer := models.RespChangePlayer{
		NextPlayer: nextPlayerId,
		Playing: false,
	}
	currResp, err := json.Marshal(currPlayer)
	if err != nil {
		logger.ZapLogger.Errorw(models.CHANGEPLAYER, "JSON Error", err)
		return
	}
	targetMap[playerId] = currResp

	logger.ZapLogger.Infoln("Exit Change Player")
	return
}

func CalculateRent(request json.RawMessage, client *models.Client, db db.DbOperations) (targetMap map[string][]byte) {
	logger.ZapLogger.Infoln("Enter Calculate Rent")	
	var req models.ReqCalculateRent
	var rent int
	targetMap = make(map[string][]byte, 2)

	err := json.Unmarshal(request, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.CALCULATERENT, "JSON Error", err)
		return
	}
	playerId := client.PlayerId
	gameId := client.GameId
	baseRent := req.Price * 10 / 100

	logger.ZapLogger.Infoln(models.CALCULATERENT, "Calculating Rent")

	switch req.BlockType {

	case models.UTILITY:

		cardCount, err := db.GetCardOwnerCount(playerId, gameId)
		if err != nil {
			logger.ZapLogger.Errorw(models.CALCULATERENT, "DB Error", err)
			return
		}
		rent = utilityRentCalculation(baseRent, cardCount)

	case models.CITY:

		status, err := db.GetCardOwnershipStatus(req.OwnerId, req.BlockId)
		if err != nil {
			logger.ZapLogger.Errorw(models.CALCULATERENT, "DB ERROR", err)
			return
		}
		rent = cityRentCalculation(baseRent, status)

	default:

	}

	playerResp := models.RespCalculateRent{
		BlockId: req.BlockId,
		BlockType: req.BlockType,
		OwnerId: req.OwnerId,
		RenterId: playerId,
		Rent: rent,
	}

	playerResponse, err := json.Marshal(playerResp)
	if err != nil {
		logger.ZapLogger.Infow(models.CALCULATERENT, "Resp for", client.PlayerId, "JSON err", err)
		return
	}

	logger.ZapLogger.Infow(models.CALCULATERENT, "Client", playerId, "Resp Body", string(playerResponse))

	targetMap[client.PlayerId] = playerResponse

	ownerResp := models.RespCalculateRent{
		BlockId: req.BlockId,
		BlockType: req.BlockType,
		RenterId: playerId,
		OwnerId: req.OwnerId,
		Rent: rent,
	}

	ownerResponse, err := json.Marshal(ownerResp)
	if err != nil {
		logger.ZapLogger.Infow(models.CALCULATERENT, "Resp for", req.OwnerId, "JSON err", err)
		return
	}
	logger.ZapLogger.Infow(models.CALCULATERENT, "Client", req.OwnerId, "Resp Body", string(ownerResponse))

	targetMap[req.OwnerId] = ownerResponse

	logger.ZapLogger.Infoln("Exit Calculate Rent")
	return
}

func PayRent(request json.RawMessage, client *models.Client, db db.DbOperations) (targetMap map[string][]byte) {
	logger.ZapLogger.Infoln("Enter Pay Rent")	
	var req models.ReqPayRent
	targetMap = make(map[string][]byte, 2)

	err := json.Unmarshal(request, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.PAYRENT, "JSON Error", err)
		return
	}
	playerId := client.PlayerId
	logger.ZapLogger.Infow(models.PAYRENT, "Owner Id", req.OwnerId)
	ownerCash, err := db.GetPlayerCash(req.OwnerId)
	if err != nil {
		logger.ZapLogger.Infow(models.PAYRENT, "DB Error", err)
		return
	}

	updatedOwnerCash := ownerCash - req.Rent
	updatedPlayerCash := req.Cash - req.Rent

	err = db.UpdatePlayerCash(playerId, updatedPlayerCash)
	if err != nil {
		logger.ZapLogger.Infow(models.PAYRENT, "DB Error", err)
		return
	}

	err = db.UpdatePlayerCash(req.OwnerId, updatedOwnerCash)
	if err != nil {
		logger.ZapLogger.Infow(models.PAYRENT, "DB Error", err)
		return
	}

	playerResp := models.RespPayRent{
		Cash: updatedPlayerCash,
		Paid: true,
		OwnerId: req.OwnerId,
	}

	PlayerResponse, err := json.Marshal(playerResp)
	if err != nil {
		logger.ZapLogger.Infow(models.PAYRENT, "Resp for", playerId, "State", models.PAYRENT, "JSON err", err)
		return
	}

	targetMap[playerId] = PlayerResponse

	ownerResp := models.RespPayRent{
		Cash: updatedOwnerCash,
		Paid: true,
		RenterId: playerId,
	}

	ownerResponse, err := json.Marshal(ownerResp)
	if err != nil {
		logger.ZapLogger.Infow(models.PAYRENT, "Resp for", req.OwnerId, "State", models.PAYRENT, "JSON err", err)
		return
	}

	targetMap[req.OwnerId] = ownerResponse

	logger.ZapLogger.Infoln("Exit Pay Rent")
	return
}

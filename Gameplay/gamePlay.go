package gameplay

import (
	db "Monopoly/DB"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
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

	go callMovePos(diceVal, client, readCh)

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

	cash, pos, err := db.GetPlayerCashPos(playerId)
	if err != nil {
		logger.ZapLogger.Errorw(models.MOVEPOS, "DB Error", err)
		return
	}
	logger.ZapLogger.Infow(models.MOVEPOS, "Current Position", pos)

	newPos := updatePos(req.UpdateBy, pos)
	// err = db.UpdatePlayerPos(playerId, newPos)
	err = db.UpdatePlayerCash(playerId, cash, newPos)
	if err != nil {
		logger.ZapLogger.Errorw(models.MOVEPOS, "DB Error", err)
	}

	logger.ZapLogger.Infow(models.MOVEPOS, "Player Id", playerId, "New Pos", newPos)
	block, err := db.GetBlockState(newPos, gameId)
	if err != nil {
		logger.ZapLogger.Errorw(models.MOVEPOS, "DB Error", err)
		return
	}

	switch block.OwnerId {
	case "":
		// Buy or Action Card
		if block.Type == models.SPECIALCARD {
			var status string
			resp = models.RespMovePos{
				BlockId:  block.BlockId,
				NewPos:   newPos,
				Type:     block.Type,
				BlockName: block.BlockName,
			}
			if block.BlockName == models.COMMUNITYCHEST || block.BlockName == models.CHANCE {
				resp.CardNo = getCardNo()
				// cardInfo, err := db.GetBlockInfo(block.BlockName, resp.CardNo)
				// if err != nil {
				// 	logger.ZapLogger.Errorw(models.MOVEPOS, "DB Error", err)
				// 	return
				// }
				// resp.CardInfo = cardInfo
				status = models.ACTIONCARD + "_" + strconv.Itoa(resp.CardNo) + "_" + block.BlockName
			}

			status = models.ACTIONCARD + "_" + block.BlockName
			err = db.UpdatePlayerStatus(playerId, status)
			if err != nil {
				logger.ZapLogger.Errorw(models.MOVEPOS, "DB Error", err)
				return
			}
			// go callActionCard(client, readCh)

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

		calRentReq := &models.ReqCalculateRent{
			BlockId: block.BlockId,
			BlockType: block.Type,
			Price: block.Price,
			OwnerId: block.OwnerId,
		}
		go callCalculateRent(calRentReq, client, readCh)
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
	buy = false

	cash, pos, err := db.GetPlayerCashPos(playerId)
	if err != nil {
		logger.ZapLogger.Errorw(models.BUYBLOCK, "DB Error", err)
		return
	}

	price, err := db.GetBlockPrice(req.BlockId)
	// block, err := db.GetBlockState(req.Pos, gameId)
	if err != nil {
		logger.ZapLogger.Errorw(models.BUYBLOCK, "DB Error", err)
		return
	}

	if checkPlayerCash(cash, price) {
		updatedCash = cash - price
		buy = true
	} else {
		logger.ZapLogger.Infow("Player Low on Cash", "Player Id", playerId, "Cash", cash)
		buy = false

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
	}

	err = db.UpdatePlayerCash(playerId, updatedCash, pos)
	if err != nil {
		logger.ZapLogger.Errorw(models.BUYBLOCK, "DB Error", err)
		return
	}
	if buy {
		go callChangePlayer(client, readCh)
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
	inJail := false

	cash, pos, err := db.GetPlayerCashPos(playerId)
	if err != nil {
		logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
		return
	}

	blockType := req.Type
	// block, err := db.GetBlockState(req.Pos, gameId)
	// if err != nil {
	// 	logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
	// 	return
	// }

	logger.ZapLogger.Infow(models.ACTIONCARD, "Current Cash", cash)

	switch blockType {

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
		cash += chestValue
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
			pos, err = db.GetPosByBlockName("Mumbai")
			if err != nil {
				logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
				return
			}

		case models.GOTOSTART:
			pos, err = db.GetPosByBlockName("Go Start")
			if err != nil {
				logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
				return
			}
			cash += 2000

		case models.HOUSEHOTELFINE:
			statusList, err := db.GetPlayerStatusList(playerId, gameId)
			if err != nil {
				logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
				return
			}

			fine := houseHotelFine(statusList)
			if checkPlayerCash(cash, fine) {
				cash -= fine
			} else {
				err = errors.New("Not Enough Cash")
				return
			}

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
			cash += value
			logger.ZapLogger.Infow(models.ACTIONCARD, "Value", value)

		}

	case models.INCOMETAX:
		price, err := db.GetBlockPrice(req.BlockId)
		if err != nil {
			logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
			return
		}
		cash -= price

	case models.JAIL:
		// resp.Pos = req.Pos
		inJail = true
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

		pos = jailPos
		inJail = true
		status := models.BLOCKED + "_3"

		err = db.UpdatePlayerStatus(playerId, status)
		if err != nil {
			logger.ZapLogger.Errorw(models.JAIL, "DB Error", err)
			return
		}

	case models.PROPERTYTAX:
		price, err := db.GetBlockPrice(req.BlockId)
		if err != nil {
			logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
			return
		}
		cash -= price

	case models.GOTOSTART:
		price, err := db.GetBlockPrice(req.BlockId)
		if err != nil {
			logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
			return
		}
		cash -= price

	}

	// Update Player Remaining
	err = db.UpdatePlayerCash(playerId, cash, pos)
	if err != nil {
		logger.ZapLogger.Infow(models.ACTIONCARD, "DB Error", err)
		return
	}

	if !inJail {
		go callChangePlayer(client, readCh)
	}

	resp := models.RespActionCard{
		Cash: cash,
		Pos: pos,
		InJail: inJail,
	}

	response, err = json.Marshal(resp)
	if err != nil {
		logger.ZapLogger.Errorw(models.ROLLDICE, "JSON Error", err)
		return
	}

	logger.ZapLogger.Infow(models.ACTIONCARD, "Resp", string(response))
	return
}

func Jail(request json.RawMessage, client *models.Client, db db.DbOperations, readCh chan<- models.WSMessage) (response json.RawMessage) {
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

	if !req.InJail {
		return
	}

	cash, pos, err := db.GetPlayerCashPos(playerId)
	if err != nil {
		logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
		return
	}

	switch req.JailId {

	case "Jail1":
		updatedCash = cash - 500
		logger.ZapLogger.Infow(models.JAIL, "Updated Cash", updatedCash)
		err = db.UpdatePlayerCash(playerId, updatedCash, pos)
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
		updatedCash = cash

		err = db.UpdatePlayerStatus(playerId, "")
		if err != nil {
			logger.ZapLogger.Errorw(models.JAIL, "DB Error", err)
			return
		}

	case "Jail3":
		// The player is blocked at the front end and they are also checked at backend
		inJail = req.InJail
		updatedCash = cash

	}
	
	resp := models.RespJail{
		Cash: updatedCash,
		InJail: inJail,
	}

	go callChangePlayer(client, readCh)

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

	status := models.PAYRENT + "-" + strconv.Itoa(rent) + "-" + req.OwnerId
	err = db.UpdatePlayerStatus(playerId, status)
	if err != nil {
		logger.ZapLogger.Errorw(models.CALCULATERENT, "DB Error", err)
		return
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

func PayRent(request json.RawMessage, client *models.Client, db db.DbOperations, readCh chan<- models.WSMessage) (targetMap map[string][]byte) {
	logger.ZapLogger.Infoln("Enter Pay Rent")	
	var req models.ReqPayRent
	targetMap = make(map[string][]byte, 2)

	err := json.Unmarshal(request, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.PAYRENT, "JSON Error", err)
		return
	}
	playerId := client.PlayerId
	gameId := client.GameId

	status, err := db.GetPlayerStatus(playerId, gameId)
	if err != nil {
		logger.ZapLogger.Errorw(models.PAYRENT, "DB Error", err)
		return
	}

	statusList := strings.Split(status, "-")
	ownerId := statusList[2]
	rent, err := strconv.Atoi(statusList[1])
	if err != nil {
		logger.ZapLogger.Errorw(models.PAYRENT, "Conversion Error", err)
		return
	}

	playerCash, playerPos, err := db.GetPlayerCashPos(playerId)
	if err != nil {
		logger.ZapLogger.Errorw(models.PAYRENT, "DB Error", err)
		return
	}

	logger.ZapLogger.Infow(models.PAYRENT, "Owner Id", ownerId)
	ownerCash, ownerPos, err := db.GetPlayerCashPos(ownerId)
	if err != nil {
		logger.ZapLogger.Infow(models.PAYRENT, "DB Error", err)
		return
	}

	updatedOwnerCash := ownerCash - rent
	updatedPlayerCash := playerCash - rent

	// Add the correct Pos
	err = db.UpdatePlayerCash(playerId, updatedPlayerCash, playerPos)
	if err != nil {
		logger.ZapLogger.Infow(models.PAYRENT, "DB Error", err)
		return
	}

	err = db.UpdatePlayerCash(ownerId, updatedOwnerCash, ownerPos)
	if err != nil {
		logger.ZapLogger.Infow(models.PAYRENT, "DB Error", err)
		return
	}

	playerResp := models.RespPayRent{
		Cash: updatedPlayerCash,
		Paid: true,
		Rent: rent,
		OwnerId: ownerId,
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
		Rent: rent,
		RenterId: playerId,
	}

	ownerResponse, err := json.Marshal(ownerResp)
	if err != nil {
		logger.ZapLogger.Infow(models.PAYRENT, "Resp for", ownerId, "State", models.PAYRENT, "JSON err", err)
		return
	}
	targetMap[ownerId] = ownerResponse

	go callChangePlayer(client, readCh)

	logger.ZapLogger.Infoln("Exit Pay Rent")
	return
}

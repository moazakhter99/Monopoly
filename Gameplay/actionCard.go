package gameplay

import (
	db "Monopoly/DB"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
	"errors"
	"strconv"
)



type ActionCardProc struct {
	db      db.DbOperations
	client  *models.Client
	writeCh chan<- models.WSMessage

}

func CreateActionCard(db db.DbOperations) *ActionCardProc {
	return &ActionCardProc{
		db: db,
	}
}

func (a *ActionCardProc) Validate(reqMsg []byte) (payload any, err error) {
	logger.ZapLogger.Infoln("Enter Validate Action Card")
	var req models.ReqActionCard
	err = json.Unmarshal(reqMsg, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.ACTIONCARD, "Validation Error", err)
		return
	}
	logger.ZapLogger.Infoln("Exit Validate Action Card")
	return req, err
}

func (a *ActionCardProc) Play(payload any) ([]byte, error) {
	logger.ZapLogger.Infoln("Enter Play Action Card")

	var err error
	playerId := a.client.PlayerId
	gameId := a.client.GameId
	inJail := false
	req := payload.(models.ReqActionCard)
	blockType := req.Type

	cash, pos, err := a.db.GetPlayerCashPos(playerId)
	if err != nil {
		logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
		return nil, err
	}

	logger.ZapLogger.Infow(models.ACTIONCARD, "Current Cash", cash)

	switch blockType {

	case models.COMMUNITYCHEST:
		action, err := a.db.GetCardAction(req.CardId)
		if err != nil {
			logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
			return nil, err
		}

		chestValue, err := strconv.Atoi(action)
		if err != nil {
			logger.ZapLogger.Errorw(models.ACTIONCARD, "Conversion Error", err)
			return nil, err
		}
		cash += chestValue
		logger.ZapLogger.Infow(models.ACTIONCARD, "Update Cash", chestValue)

	case models.CHANCE:
	
		action, err := a.db.GetCardAction(req.CardId)
		if err != nil {
			logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
			return nil, err
		}

		logger.ZapLogger.Infow(models.ACTIONCARD, "Action", action)
		switch action {

		case models.JUMPTOMUMBAI:
			pos, err = a.db.GetPosByBlockName("Mumbai")
			if err != nil {
				logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
				return nil, err
			}

		case models.GOTOSTART:
			pos, err = a.db.GetPosByBlockName("Go Start")
			if err != nil {
				logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
				return nil, err
			}
			cash += 2000

		case models.HOUSEHOTELFINE:
			statusList, err := a.db.GetPlayerStatusList(playerId, gameId)
			if err != nil {
				logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
				return nil, err
			}

			fine := houseHotelFine(statusList)
			if checkPlayerCash(cash, fine) {
				cash -= fine
			} else {
				err = errors.New("Not Enough Cash")
				return nil, err
			}

		case models.GETOUTOFJAIL:
			err = a.db.UpdateGetOutOfJailCard(playerId, gameId)
			if err != nil {
				logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
				return nil, err
			}

		default:
			value, err := strconv.Atoi(action)
			if err != nil {
				logger.ZapLogger.Infow(models.ACTIONCARD, "Conversion Error", err, "Invalid Action", action)
				return nil, err
			}
			cash += value
			logger.ZapLogger.Infow(models.ACTIONCARD, "Value", value)

		}

	case models.INCOMETAX:
		price, err := a.db.GetBlockPrice(req.BlockId)
		if err != nil {
			logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
			return nil, err
		}
		cash -= price

	case models.JAIL:
		// resp.Pos = req.Pos
		inJail = true
		status := models.BLOCKED + "_3"

		err = a.db.UpdatePlayerStatus(playerId, status)
		if err != nil {
			logger.ZapLogger.Errorw(models.JAIL, "DB Error", err)
			return nil, err
		}

	case models.FREEPARKING:

	case models.GOTOJAIL:

		jailPos, err := a.db.GetPosByBlockName("Jail")
		if err != nil {
			logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
			return nil, err
		}

		pos = jailPos
		inJail = true
		status := models.BLOCKED + "_3"

		err = a.db.UpdatePlayerStatus(playerId, status)
		if err != nil {
			logger.ZapLogger.Errorw(models.JAIL, "DB Error", err)
			return nil, err
		}

	case models.PROPERTYTAX:
		price, err := a.db.GetBlockPrice(req.BlockId)
		if err != nil {
			logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
			return nil, err
		}
		cash -= price

	case models.GOTOSTART:
		price, err := a.db.GetBlockPrice(req.BlockId)
		if err != nil {
			logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
			return nil, err
		}
		cash -= price

	}

	// Update Player Remaining
	err = a.db.UpdatePlayerCash(playerId, cash, pos)
	if err != nil {
		logger.ZapLogger.Infow(models.ACTIONCARD, "DB Error", err)
		return nil, err
	}

	// if !inJail {
	// 	go callChangePlayer(client, readCh)
	// }

	response := models.RespActionCard{
		Cash: cash,
		Pos: pos,
		InJail: inJail,
	}

	resp, err := json.Marshal(response)
	if err != nil {
		logger.ZapLogger.Errorw(models.ROLLDICE, "JSON Error", err)
		return nil, err
	}

	logger.ZapLogger.Infoln("Exit Play Action Card")
	return resp, nil
}

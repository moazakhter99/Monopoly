package gameplay

import (
	db "Monopoly/DB"
	gameroom "Monopoly/Gameroom"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
	"strconv"
	"strings"
)

type PayRentProc struct {
	db   db.DbOperations
	room gameroom.Room
}

func CreatePayRent(db db.DbOperations, room gameroom.Room) *PayRentProc {
	return &PayRentProc{
		db:   db,
		room: room,
	}
}

func (p *PayRentProc) Validate(reqMsg []byte) (payload any, err error) {
	logger.ZapLogger.Infoln("Enter Validate Pay Rent")
	var req models.ReqPayRent
	err = json.Unmarshal(reqMsg, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.PAYRENT, "Validation Error", err)
		return
	}
	logger.ZapLogger.Infoln("Exit Validate Pay Rent")
	return req, err
}

func (p *PayRentProc) Play(payload any, param map[string]string) (targetMap map[string]any, err error) {
	logger.ZapLogger.Infoln("Enter Play Pay Rent")
	// req := payload.(models.ReqPayRent)

	gameId := param["Game"]
	playerId := param["Player"]
	targetMap = make(map[string]any, 2)

	status, err := p.db.GetPlayerStatus(playerId, gameId)
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

	playerCash, playerPos, err := p.db.GetPlayerCashPos(playerId)
	if err != nil {
		logger.ZapLogger.Errorw(models.PAYRENT, "DB Error", err)
		return
	}

	logger.ZapLogger.Infow(models.PAYRENT, "Owner Id", ownerId)
	ownerCash, ownerPos, err := p.db.GetPlayerCashPos(ownerId)
	if err != nil {
		logger.ZapLogger.Infow(models.PAYRENT, "DB Error", err)
		return
	}

	updatedOwnerCash := ownerCash - rent
	updatedPlayerCash := playerCash - rent

	// Add the correct Pos
	err = p.db.UpdatePlayerCash(playerId, updatedPlayerCash, playerPos)
	if err != nil {
		logger.ZapLogger.Infow(models.PAYRENT, "DB Error", err)
		return
	}

	err = p.db.UpdatePlayerCash(ownerId, updatedOwnerCash, ownerPos)
	if err != nil {
		logger.ZapLogger.Infow(models.PAYRENT, "DB Error", err)
		return
	}

	playerResp := models.RespPayRent{
		Cash:    updatedPlayerCash,
		Paid:    true,
		Rent:    rent,
		OwnerId: ownerId,
	}
	targetMap[playerId] = playerResp

	ownerResp := models.RespPayRent{
		Cash:     updatedOwnerCash,
		Paid:     true,
		Rent:     rent,
		RenterId: playerId,
	}
	targetMap[ownerId] = ownerResp

	p.room.UpdateGameState(gameId, playerId, models.PAYRENT)

	logger.ZapLogger.Infoln("Exit Play Pay Rent")
	return
}

// Response implements [Game].
func (p *PayRentProc) Response(targetMap map[string]any, reqParam map[string]string, readChan chan []byte) (err error) {
	// go callChangePlayer(client, readCh)
	panic("unimplemented")
}

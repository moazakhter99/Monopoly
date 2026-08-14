package gameplay

import (
	db "Monopoly/DB"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
	"strconv"
	"strings"
)



type PayRentProc struct {
	db 		db.DbOperations
	client 	*models.Client
	writeCh chan<- models.WSMessage

}

func CreatePayRent(db db.DbOperations) *PayRentProc {
	return &PayRentProc{
		db: db,
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

func (p *PayRentProc) Play(payload any) (resp []byte, err error) {
	logger.ZapLogger.Infoln("Enter Play Pay Rent")
	// req := payload.(models.ReqPayRent)
	targetMap := make(map[string][]byte, 2)

	playerId := p.client.PlayerId
	gameId := p.client.GameId

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

	// go callChangePlayer(client, readCh)

	logger.ZapLogger.Infoln("Exit Play Pay Rent")
	return
}

package gameplay

import (
	db "Monopoly/DB"
	gameroom "Monopoly/Gameroom"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
	"strconv"
)

type CalculateRentProc struct {
	db   db.DbOperations
	room gameroom.Room
}

func CreateCalculateRent(db db.DbOperations, room gameroom.Room) *CalculateRentProc {
	return &CalculateRentProc{
		db:   db,
		room: room,
	}
}

func (c *CalculateRentProc) Validate(reqMsg []byte) (payload any, err error) {
	logger.ZapLogger.Infoln("Enter Validate Calculate Rent")
	var req models.ReqCalculateRent
	err = json.Unmarshal(reqMsg, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.CALCULATERENT, "Validation Error", err)
		return
	}
	logger.ZapLogger.Infoln("Exit Validate Calculate Rent")
	return req, err
}

// Might not work as the resposens are targeted
func (c *CalculateRentProc) Play(payload any, param map[string]string) (map[string]any, error) {
	logger.ZapLogger.Infoln("Enter Play Calculate Rent")
	req := payload.(models.ReqCalculateRent)
	var rent int
	var err error
	targetMap := make(map[string]any, 2)

	gameId := param["Game"]
	playerId := param["Player"]
	baseRent := req.Price * 10 / 100

	logger.ZapLogger.Infoln(models.CALCULATERENT, "Calculating Rent")

	switch req.BlockType {

	case models.UTILITY:

		cardCount, err := c.db.GetCardOwnerCount(playerId, gameId)
		if err != nil {
			logger.ZapLogger.Errorw(models.CALCULATERENT, "DB Error", err)
			return nil, err
		}
		rent = utilityRentCalculation(baseRent, cardCount)

	case models.CITY:

		status, err := c.db.GetCardOwnershipStatus(req.OwnerId, req.BlockId)
		if err != nil {
			logger.ZapLogger.Errorw(models.CALCULATERENT, "DB ERROR", err)
			return nil, err
		}
		rent = cityRentCalculation(baseRent, status)

	default:

	}

	status := models.PAYRENT + "-" + strconv.Itoa(rent) + "-" + req.OwnerId
	err = c.db.UpdatePlayerStatus(playerId, status)
	if err != nil {
		logger.ZapLogger.Errorw(models.CALCULATERENT, "DB Error", err)
		return nil, err
	}

	playerResp := models.RespCalculateRent{
		BlockId:   req.BlockId,
		BlockType: req.BlockType,
		OwnerId:   req.OwnerId,
		RenterId:  playerId,
		Rent:      rent,
	}

	logger.ZapLogger.Infow(models.CALCULATERENT, "Client", playerId, "Resp Body", playerResp)

	targetMap[playerId] = playerResp

	ownerResp := models.RespCalculateRent{
		BlockId:   req.BlockId,
		BlockType: req.BlockType,
		RenterId:  playerId,
		OwnerId:   req.OwnerId,
		Rent:      rent,
	}

	logger.ZapLogger.Infow(models.CALCULATERENT, "Client", req.OwnerId, "Resp Body", ownerResp)

	targetMap[req.OwnerId] = ownerResp

	logger.ZapLogger.Infoln("Exit Play Calculate Rent")
	return targetMap, nil
}

// Response implements [Game].
func (c *CalculateRentProc) Response(targetMap map[string]any, reqParam map[string]string, readChan chan []byte) (err error) {
	panic("unimplemented")
}
